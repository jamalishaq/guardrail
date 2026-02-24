// tool/tool.go
package tool

import "context"

// Tool is a capability exposed by Guardrail.
// Behavior-first: identity is stable (Name), description is separate (Spec).
type Tool interface {
	// Name is the stable identifier used for:
	// - registry keys
	// - tool calls (MCP ToolCall.name)
	// - audit keys, rate-limit keys, policy lookup, etc.
	Name() string

	// Spec provides metadata used for discovery/listing and input validation schema.
	Spec() Spec

	// Execute runs the tool.
	Execute(ctx context.Context, input any) (any, error)
}