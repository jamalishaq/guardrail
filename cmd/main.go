// cmd/demo/main.go
// Demo showing a tool that returns a Spec with JSON schema similar to MCP tools/list.
// go run ./cmd/demo
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jamalishaq/guardrail/execution"
	"github.com/jamalishaq/guardrail/tool"
)

type calculatorArithmeticTool struct{}

func (t calculatorArithmeticTool) Name() string { return "calculator_arithmetic" }

func (t calculatorArithmeticTool) Spec() tool.Spec {
	return tool.Spec{
		Title:       "Calculator",
		Description: "Perform mathematical calculations including basic arithmetic, trigonometric functions, and algebraic operations",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"expression": map[string]any{
					"type":        "string",
					"description": "Mathematical expression to evaluate (e.g., '2 + 3 * 4', 'sin(30)', 'sqrt(16)')",
				},
			},
			"required": []any{"expression"},
		},
	}
}

func (t calculatorArithmeticTool) Execute(ctx context.Context, input any) (any, error) {
	// MVP: pretend we computed something. We'll implement real parsing later.
	m, _ := input.(map[string]any)
	expr, _ := m["expression"].(string)
	if expr == "" {
		expr = "(empty)"
	}
	return map[string]any{
		"expression": expr,
		"result":     "not_implemented_yet",
	}, nil
}

func main() {
	reg := tool.NewRegistry()
	_ = reg.Add(calculatorArithmeticTool{})

	exec := execution.NewExecutor(reg)

	// Show what would be used for MCP tools/list
	for _, t := range reg.ListTools() {
		fmt.Printf("Tool: %s | title=%q\n", t.Name(), t.Spec().Title)
	}

	inv := execution.Invocation{
		ID:        "inv-1",
		ToolName:  "calculator_arithmetic",
		Input:     map[string]any{"expression": "2 + 3 * 4"},
		Caller:    execution.Caller{ID: "user-123", Type: "user"},
		Metadata:  execution.Metadata{"request_id": "req-abc"},
		StartTime: time.Now(),
	}

	res := exec.Execute(context.Background(), inv)
	fmt.Println("Output:", res.Output)
	fmt.Println("Error:", res.Error)
	fmt.Println("Duration:", res.Duration)
}