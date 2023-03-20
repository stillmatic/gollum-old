package gollum

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/antonmedv/expr"
)

type Tool struct {
	// Name is the name of the tool, will be used for lookup
	Name string
	// Description is a short description of the tool with usage info
	Description string
	// Run is the function that will be called when the tool is invoked
	Run func(arg string) (string, error)
}

// WikipediaResult is the result of a Wikipedia search
type WikipediaResult struct {
	Batchcomplete string `json:"batchcomplete"`
	Continue      struct {
		Sroffset int    `json:"sroffset"`
		Continue string `json:"continue"`
	} `json:"continue"`
	Query struct {
		Searchinfo struct {
			Totalhits         int    `json:"totalhits"`
			Suggestion        string `json:"suggestion"`
			Suggestionsnippet string `json:"suggestionsnippet"`
		} `json:"searchinfo"`
		Search []struct {
			Ns        int       `json:"ns"`
			Title     string    `json:"title"`
			Pageid    int       `json:"pageid"`
			Size      int       `json:"size"`
			Wordcount int       `json:"wordcount"`
			Snippet   string    `json:"snippet"`
			Timestamp time.Time `json:"timestamp"`
		} `json:"search"`
	} `json:"query"`
}

// RunWikipedia queries Wikipedia for the given search term and returns the
// snippet of the first result
func RunWikipedia(arg string) (string, error) {
	var result WikipediaResult
	queryParams := url.Values{
		"action":   []string{"query"},
		"list":     []string{"search"},
		"srsearch": []string{arg},
		"format":   []string{"json"},
		"srliimit": []string{"1"},
	}
	resp, err := http.Get("http://en.wikipedia.org/w/api.php?" + queryParams.Encode())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if len(result.Query.Search) == 0 {
		return "", errors.New("no results")
	}
	snippet := result.Query.Search[0].Snippet
	snippet = strings.Replace(snippet, "<span class=\"searchmatch\">", "", -1)
	snippet = strings.Replace(snippet, "</span>", "", -1)
	return snippet, nil
}

// RunCalculator evaluates  mathematical expressions and returns the result.
// Internally, this uses the `expr` package to avoid arbitrary code execution.
func RunCalculator(arg string) (string, error) {
	env := map[string]interface{}{}
	program, err := expr.Compile(arg, expr.Env(env))
	if err != nil {
		return "", err
	}
	output, err := expr.Run(program, nil)
	if err != nil {
		return "", err
	}
	switch t := output.(type) {
	case string:
		return t, nil
	case int:
		return strconv.Itoa(t), nil
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64), nil
	default:
		return "", errors.New("invalid output")
	}
}

var WikipediaTool = Tool{
	Name:        "wikipedia",
	Description: "Search Wikipedia for the given term and return the first result",
	Run:         RunWikipedia,
}

var CalculatorTool = Tool{
	Name:        "calculator",
	Description: "Evaluate mathematical expressions",
	Run:         RunCalculator,
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
	return tool.Run(strings.TrimSpace(parts[2]))
}
