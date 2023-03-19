package react_test

import (
	"context"
	"os"
	"testing"

	"github.com/PullRequestInc/go-gpt3"
	react "github.com/stillmatic/go-llm-react"
	"github.com/stretchr/testify/assert"
)

const (
	testConvoName = "test"
)

func TestReactEndToEnd(t *testing.T) {
	reg := react.NewToolRegistry()
	apiKey := os.Getenv("OPENAI_KEY")
	if apiKey == "" {
		t.Fatal("OPENAI_KEY is not set")
	}
	aiClient := gpt3.NewClient(apiKey)
	r := react.NewReactAgent(aiClient, reg)
	r.NewConversation(testConvoName)
	ctx := context.Background()
	err := r.Speak(ctx, testConvoName, "Question: What does England share borders with?")
	assert.NoError(t, err)
}
