package gollum

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/PullRequestInc/go-gpt3"
)

const (
	RoleSystem    = "system"
	RoleUser      = "user"
	RoleAssistant = "assistant"
)

type ReactAgent struct {
	Client        gpt3.Client
	Registry      *ToolRegistry
	Conversations map[string]*Conversation
	MaxTurns      int
}

type Conversation struct {
	Messages  []gpt3.ChatCompletionRequestMessage
	CurrReply *strings.Builder
}

//go:embed prompt.txt
var initialPrompt string

func NewReactAgent(client gpt3.Client, registry *ToolRegistry) *ReactAgent {
	return &ReactAgent{
		Client:        client,
		Registry:      registry,
		Conversations: make(map[string]*Conversation),
		MaxTurns:      10,
	}
}

func (a *ReactAgent) NewConversation(name string) {
	conv := &Conversation{
		Messages: []gpt3.ChatCompletionRequestMessage{
			{
				Role:    RoleSystem,
				Content: initialPrompt,
			},
		},
		CurrReply: &strings.Builder{},
	}
	a.Conversations[name] = conv
}

func (a *ReactAgent) Speak(ctx context.Context, conversationName string, prompt string) error {
	conv, ok := a.Conversations[conversationName]
	if !ok {
		return errors.New("conversation not found")
	}
	conv.Messages = append(conv.Messages, gpt3.ChatCompletionRequestMessage{
		Role:    RoleUser,
		Content: prompt,
	})
	for {
		if len(conv.Messages) >= a.MaxTurns {
			fmt.Println("Max turns reached")
			break
		}
		done, err := a.speak(ctx, conv)
		if err != nil {
			return err
		}
		if done {
			break
		}
	}
	return nil
}

func (c *Conversation) onData(data *gpt3.ChatCompletionStreamResponse) {
	msg := data.Choices[0].Delta.Content
	c.CurrReply.WriteString(msg)
	fmt.Print(msg)
}

func (a *ReactAgent) speak(ctx context.Context, conv *Conversation) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err := a.Client.ChatCompletionStream(ctx, gpt3.ChatCompletionRequest{
		Model:       gpt3.GPT3Dot5Turbo,
		Messages:    conv.Messages,
		MaxTokens:   256,
		Temperature: 0,
		Stop:        []string{"PAUSE"},
	}, conv.onData)
	if err != nil {
		return false, err
	}
	respMessage := conv.CurrReply.String()
	conv.CurrReply.Reset()
	conv.Messages = append(conv.Messages, gpt3.ChatCompletionRequestMessage{
		Role:    RoleAssistant,
		Content: respMessage,
	})
	obs, err := a.Registry.Run(respMessage)
	if err != nil {
		if errors.Is(err, ErrNoActionFound) {
			return true, nil
		}
		return false, err
	}
	fmt.Println("\nObservation: " + obs)
	conv.Messages = append(conv.Messages, gpt3.ChatCompletionRequestMessage{
		Role:    RoleSystem,
		Content: "Observation: " + obs,
	})
	return false, nil
}
