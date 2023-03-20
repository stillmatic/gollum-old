package tools

import (
	"errors"
	"strconv"

	"github.com/antonmedv/expr"
)

var CalculatorTool = Tool{
	Name:        "calculator",
	Description: "Evaluate mathematical expressions",
	Run:         RunCalculator,
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
