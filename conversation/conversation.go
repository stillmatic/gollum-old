package conversation

import (
	"fmt"
	"strings"

	"github.com/PullRequestInc/go-gpt3"
)

type Conversation struct {
	Messages  []gpt3.ChatCompletionRequestMessage
	CurrReply *strings.Builder
}

func (c *Conversation) OnDataPrint(data *gpt3.ChatCompletionStreamResponse) {
	msg := data.Choices[0].Delta.Content
	c.CurrReply.WriteString(msg)
	fmt.Print(msg)
}
