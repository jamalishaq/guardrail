// execution/invocation.go
package execution

import "time"

// Metadata is lightweight key/value context attached to an invocation.
// Keep values safe-to-log by default (no secrets).
type Metadata map[string]string

// Caller represents the identity of the invoker (user/service/agent).
// For the MVP, keep it simple. You can expand later (scopes, auth method, etc.).
type Caller struct {
	ID       string   // stable principal ID
	Type     string   // e.g. "user", "service", "agent"
	Roles    []string // e.g. ["admin", "sales"]
	TenantID string   // optional
}

// Invocation is Guardrail's internal representation of a tool call.
// It is protocol-agnostic and should NOT contain MCP SDK types.
type Invocation struct {
	ID        string
	ToolName  string
	Input     any
	Caller    Caller
	Metadata  Metadata
	StartTime time.Time
}