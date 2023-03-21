package tools_test

import (
	"strings"
	"testing"

	"github.com/stillmatic/gollum/tools"
	"github.com/stretchr/testify/assert"
)

func setupRegistry(t *testing.T) *tools.ToolRegistry {
	t.Helper()
	reg := tools.NewToolRegistry()
	reg.Register(tools.CalculatorTool)
	reg.Register(tools.WikipediaTool)
	return reg
}

func TestReact(t *testing.T) {
	reg := setupRegistry(t)
	resp, err := reg.Run(`Question: What does England share borders with?
Thought: I should list down the neighboring countries of England
Action: wikipedia: England
PAUSE`)
	assert.NoError(t, err)
	assert.Equal(t, `England is a country that is part of the United Kingdom. It shares land borders with Wales to its west and Scotland to its north. The Irish Sea lies northwest and the Celtic Sea area of the Atlantic Ocean to the southwest. It is separated from continental Europe by the North Sea to the east and the English Channel to the south. The country covers five-eighths of the island of Great Britain, which lies in the North Atlantic, and includes over 100 smaller islands, such as the Isles of Scilly and the Isle of Wight.
The area now called England was first inhabited by modern humans during the Upper Paleolithic period, but takes its name from the Angles, a Germanic tribe deriving its name from the Anglia peninsula, who settled during the 5th and 6th centuries.`, resp)

	resp, err = reg.Run(`Question: What is (6 * 70) / 5?
Thought: I should use a calculator
Action: calc: (6 * 70) / 5
PAUSE`)
	assert.NoError(t, err)
	assert.Equal(t, `84`, resp)
}

func TestAddTool(t *testing.T) {
	reg := setupRegistry(t)
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
