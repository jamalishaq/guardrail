// tool/spec.go
package tool

// Spec describes a tool for discovery/listing (e.g., MCP tools/list).
// Keep it protocol-agnostic and safe-to-log (no secrets).
type Spec struct {
	Title       string
	Description string

	// InputSchema is JSON Schema represented as a generic map for the MVP.
	// Later you can switch to a typed schema representation or a JSON raw message.
	InputSchema map[string]any
}