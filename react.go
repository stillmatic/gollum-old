package react

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
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
	Messages []gpt3.ChatCompletionRequestMessage
}

//go:embed prompt.txt
var initialPrompt string

func NewReactAgent(client gpt3.Client, registry *ToolRegistry) *ReactAgent {
	return &ReactAgent{
		Client:        client,
		Registry:      registry,
		Conversations: make(map[string]*Conversation),
		MaxTurns:      6,
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

func (a *ReactAgent) speak(ctx context.Context, conv *Conversation) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	resp, err := a.Client.ChatCompletion(ctx, gpt3.ChatCompletionRequest{
		Model:     gpt3.GPT3Dot5Turbo,
		Messages:  conv.Messages,
		MaxTokens: 256,
		Stop:      []string{"PAUSE"},
	})
	if err != nil {
		return false, err
	}
	respMessage := resp.Choices[0].Message.Content
	fmt.Print(respMessage)
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
	fmt.Print("Observation: " + obs)
	conv.Messages = append(conv.Messages, gpt3.ChatCompletionRequestMessage{
		Role:    RoleSystem,
		Content: "Observation: " + obs,
	})
	return false, nil
}
