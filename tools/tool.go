package tools

import (
	"errors"
	"strings"
)

type Tool struct {
	// Name is the name of the tool, will be used for lookup
	Name string
	// Description is a short description of the tool with usage info
	Description string
	// Run is the function that will be called when the tool is invoked
	Run func(arg string) (string, error)
}

type ToolRegistry struct {
	tools map[string]Tool
}

func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: map[string]Tool{
			"wikipedia":  WikipediaTool,
			"calculator": CalculatorTool,
		},
	}
}

func (r *ToolRegistry) Register(tool Tool) {
	r.tools[strings.ToLower(tool.Name)] = tool
}

func (r *ToolRegistry) AvailableTools() string {
	tools := make([]string, len(r.tools))
	i := 0
	for tool, _ := range r.tools {
		tools[i] = tool
		i++
	}
	return strings.Join(tools, ", ")
}

var (
	ErrToolNotFound  = errors.New("tool not found")
	ErrNoActionFound = errors.New("no Action command found")
	ErrInvalidAction = errors.New("invalid Action command")
)

// Run finds the last line in the given string starting with "Action",
// extracts the tool name and runs the tool with the rest of the line as
// argument.
func (r *ToolRegistry) Run(arg string) (string, error) {
	lines := strings.Split(arg, "\n")
	var line string
	for i := len(lines) - 1; i >= 0; i-- {
		currLine := strings.TrimSpace(lines[i])
		if strings.HasPrefix(currLine, "Action") {
			line = currLine
			break
		}
	}
	if line == "" {
		return "", ErrNoActionFound
	}
	parts := strings.SplitN(strings.TrimSpace(line), " ", 3)
	if len(parts) < 3 {
		return "", ErrInvalidAction
	}
	tool, ok := r.tools[strings.Trim(strings.ToLower(strings.TrimSpace(parts[1])), ":")]
	if !ok {
		return "", ErrToolNotFound
	}
	// call tool only with parts after the tool name
	return tool.Run(strings.TrimSpace(parts[2]))
}
