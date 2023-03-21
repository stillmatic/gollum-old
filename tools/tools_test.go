package tools_test

import (
	"strings"
	"testing"

	"github.com/stillmatic/gollum/tools"
	"github.com/stretchr/testify/assert"
)

func TestReact(t *testing.T) {
	reg := tools.NewToolRegistry()
	resp, err := reg.Run(`Question: What does England share borders with?
Thought: I should list down the neighboring countries of England
Action: wikipedia: England
PAUSE`)
	assert.NoError(t, err)
	assert.Equal(t, `<span class="searchmatch">England</span> is a country that is part of the United Kingdom. It shares land borders with Wales to its west and Scotland to its north. The Irish Sea lies northwest`, resp)

	resp, err = reg.Run(`Question: What is (6 * 70) / 5?
Thought: I should use a calculator
Action: calculator: (6 * 70) / 5
PAUSE`)
	assert.NoError(t, err)
	assert.Equal(t, `84`, resp)
}

func TestAddTool(t *testing.T) {
	reg := tools.NewToolRegistry()
	assert.Equal(t, 2, len(strings.Split(reg.AvailableTools(), ",")))
	reg.Register(tools.Tool{
		Name:        "test",
		Description: "test",
		Run: func(arg string) (string, error) {
			return arg, nil
		},
	})
	assert.Equal(t, 3, len(strings.Split(reg.AvailableTools(), ",")))
}
