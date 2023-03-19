package react_test

import (
	"testing"

	react "github.com/stillmatic/go-llm-react"
	"github.com/stretchr/testify/assert"
)

func TestReact(t *testing.T) {
	reg := react.NewToolRegistry()
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
