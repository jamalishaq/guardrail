// execution/executor.go
package execution

import (
	"context"
	"errors"
	"time"

	"github.com/jamalishaq/guardrail/tool"
)

// Executor looks up tools and executes them.
// Policies/middleware will be added later; keep the MVP minimal.
type Executor struct {
	registry *tool.Registry
}

func NewExecutor(registry *tool.Registry) *Executor {
	return &Executor{registry: registry}
}

var ErrNilRegistry = errors.New("registry is nil")

// Execute runs the tool specified by inv.ToolName using inv.Input.
func (e *Executor) Execute(ctx context.Context, inv Invocation) Result {
	start := time.Now()

	if e == nil || e.registry == nil {
		return Result{Error: ErrNilRegistry, Duration: time.Since(start)}
	}
	if inv.ToolName == "" {
		return Result{Error: errors.New("invocation tool name is empty"), Duration: time.Since(start)}
	}

	t, err := e.registry.Get(inv.ToolName)
	if err != nil {
		return Result{Error: err, Duration: time.Since(start)}
	}

	out, err := t.Execute(ctx, inv.Input)
	return Result{
		Output:   out,
		Error:    err,
		Duration: time.Since(start),
	}
}