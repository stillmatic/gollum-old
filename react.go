package gollum

import (
	"context"
	_ "embed"
	"fmt"
	"strings"
	"time"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/pkg/errors"
	"github.com/stillmatic/gollum/conversation"
	"github.com/stillmatic/gollum/store"
	"github.com/stillmatic/gollum/tools"
)

const (
	RoleSystem    = "system"
	RoleUser      = "user"
	RoleAssistant = "assistant"
)

//go:embed prompt.txt
var initialPrompt string

type ReactAgent struct {
	Client        gpt3.Client
	Registry      *tools.ToolRegistry
	Conversations store.Store
	Temperature   float64
	MaxTurns      int
	MaxTokens     int
	Model         string
}

type ReactAgentOption func(*ReactAgent)

func WithMaxTurns(maxTurns int) ReactAgentOption {
	return func(a *ReactAgent) {
		a.MaxTurns = maxTurns
	}
}

func WithTools(tools ...tools.Tool) ReactAgentOption {
	return func(a *ReactAgent) {
		for _, tool := range tools {
			a.Registry.Register(tool)
		}
	}
}

func WithTemperature(temperature float64) ReactAgentOption {
	return func(a *ReactAgent) {
		a.Temperature = temperature
	}
}

func WithMaxTokens(maxTokens int) ReactAgentOption {
	return func(a *ReactAgent) {
		a.MaxTokens = maxTokens
	}
}

func WithModel(model string) ReactAgentOption {
	return func(a *ReactAgent) {
		a.Model = model
	}
}

func WithStore(store store.Store) ReactAgentOption {
	return func(a *ReactAgent) {
		a.Conversations = store
	}
}

// NewReactAgent creates a new ReactAgent with the given client and tool registry.
func NewReactAgent(client gpt3.Client, registry *tools.ToolRegistry) *ReactAgent {
	return &ReactAgent{
		Client:        client,
		Registry:      registry,
		Conversations: store.NewMemoryStore(),
		MaxTurns:      10,
		Temperature:   0.0,
		MaxTokens:     256,
		Model:         gpt3.GPT3Dot5Turbo,
	}
}

// NewReactAgentWithOpts creates a new ReactAgent with the given client and options.
func NewReactAgentWithOpts(client gpt3.Client, opts ...ReactAgentOption) *ReactAgent {
	agent := NewReactAgent(client, tools.NewToolRegistry())
	for _, opt := range opts {
		opt(agent)
	}
	return agent
}

func (a *ReactAgent) GetPrompt() string {
	sb := &strings.Builder{}
	sb.WriteString(initialPrompt)
	sb.WriteString(a.Registry.GetPrompt())
	return sb.String()
}

func (a *ReactAgent) NewConversation(name string) {
	conv := &conversation.Conversation{
		Messages: []gpt3.ChatCompletionRequestMessage{
			{
				Role:    RoleSystem,
				Content: a.GetPrompt(),
			},
		},
		CurrReply: &strings.Builder{},
	}
	a.Conversations.Set(name, conv)
}

func (a *ReactAgent) Speak(ctx context.Context, conversationName string, prompt string) error {
	conv, err := a.Conversations.Get(conversationName)
	if err != nil {
		return errors.Wrap(err, "failed to get conversation")
	}
	// if len(conv.Messages) == 1 {
	// 	fmt.Println(conv.Messages[0].Content)
	// }
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

func (a *ReactAgent) speak(ctx context.Context, conv *conversation.Conversation) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	err := a.Client.ChatCompletionStream(ctx, gpt3.ChatCompletionRequest{
		Model:       a.Model,
		Messages:    conv.Messages,
		MaxTokens:   a.MaxTokens,
		Temperature: float32(a.Temperature),
		Stop:        []string{"PAUSE"},
	}, conv.OnDataPrint)
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
	var nextMessage string
	if err != nil {
		if errors.Is(err, tools.ErrNoActionFound) {
			return true, nil
		}
		nextMessage = "\nOops: " + err.Error() + ", available tools are: " + a.Registry.AvailableTools()
	} else {
		nextMessage = "\nObservation: " + obs
	}
	fmt.Println(nextMessage)
	conv.Messages = append(conv.Messages, gpt3.ChatCompletionRequestMessage{
		Role:    RoleSystem,
		Content: nextMessage,
	})
	return false, nil
}
