package mcpadapter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/jamalishaq/guardrail/execution"
	"github.com/jamalishaq/guardrail/tool"
)

var ErrNilDeps = errors.New("adapter: nil dependency")

// Adapter wires Guardrail (registry + executor) to an MCP server.
type Adapter struct {
	server   *mcp.Server
	registry *tool.Registry
	executor *execution.Executor

	idFn func() string
}

// New creates a new MCP adapter and underlying MCP server.
// implName/implVersion identify your server to MCP clients.
func New(implName, implVersion string, reg *tool.Registry, exec *execution.Executor) (*Adapter, error) {
	if reg == nil || exec == nil {
		return nil, ErrNilDeps
	}
	srv := mcp.NewServer(&mcp.Implementation{Name: implName, Version: implVersion}, nil)

	a := &Adapter{
		server:   srv,
		registry: reg,
		executor: exec,
		idFn: func() string {
			return fmt.Sprintf("inv-%d", time.Now().UnixNano())
		},
	}
	return a, nil
}

// Server exposes the underlying MCP server (useful if you want to add resources/prompts later).
func (a *Adapter) Server() *mcp.Server { return a.server }

// RegisterAllTools registers every tool in the registry onto the MCP server.
func (a *Adapter) RegisterAllTools() error {
	if a == nil || a.server == nil || a.registry == nil {
		return ErrNilDeps
	}

	for _, t := range a.registry.ListTools() {
		if err := a.registerTool(t); err != nil {
			return err
		}
	}
	return nil
}

func (a *Adapter) registerTool(t tool.Tool) error {
	name := t.Name()
	if name == "" {
		return errors.New("adapter: tool name is empty")
	}

	spec := t.Spec()

	// MCP expects Tool.InputSchema to JSON-marshal to an object schema. :contentReference[oaicite:1]{index=1}
	if spec.InputSchema == nil {
		spec.InputSchema = map[string]any{"type": "object"}
	}

	schemaBytes, err := json.Marshal(spec.InputSchema)
	if err != nil {
		return fmt.Errorf("adapter: marshal input schema for %q: %w", name, err)
	}

	mcpTool := &mcp.Tool{
		Name:        name,
		Title:       spec.Title,
		Description: spec.Description,
		// json.RawMessage is allowed for server-provided schema values. :contentReference[oaicite:2]{index=2}
		InputSchema: json.RawMessage(schemaBytes),
	}

	// Use low-level handler so *we* control:
	// - decoding raw args
	// - validation (later)
	// - shaping CallToolResult (Content / StructuredContent / IsError)
	// Server.AddTool + ToolHandler are the right primitives for that. :contentReference[oaicite:3]{index=3}
	a.server.AddTool(mcpTool, a.makeHandler(name))
	return nil
}

func (a *Adapter) makeHandler(toolName string) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// IMPORTANT:
		// Returning error from ToolHandler is treated as a protocol error. :contentReference[oaicite:4]{index=4}
		// Tool execution failures should come back as CallToolResult with IsError=true.

		// Decode raw arguments.
		// For Server.AddTool, req.Params.Arguments is raw JSON. :contentReference[oaicite:5]{index=5}
		var args map[string]any
		if raw := req.Params.Arguments; len(raw) > 0 {
			if err := json.Unmarshal(raw, &args); err != nil {
				return &mcp.CallToolResult{
					IsError: true,
					Content: []mcp.Content{
						&mcp.TextContent{Text: "invalid tool arguments: " + err.Error()},
					},
				}, nil
			}
		} else {
			args = map[string]any{}
		}

		inv := execution.Invocation{
			ID:        a.idFn(),
			ToolName:  toolName,
			Input:     args,
			StartTime: time.Now(),
			// Caller/Metadata can be added once you attach auth/transport metadata into ctx.
		}

		res := a.executor.Execute(ctx, inv)

		if res.Error != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					&mcp.TextContent{Text: res.Error.Error()},
				},
			}, nil
		}

		// MVP: return BOTH:
		// - StructuredContent (if it marshals)
		// - Content as JSON text (always safe for LLMs)
		outText := stringify(res.Output)

		return &mcp.CallToolResult{
			StructuredContent: res.Output,
			Content: []mcp.Content{
				&mcp.TextContent{Text: outText},
			},
		}, nil
	}
}

func stringify(v any) string {
	if v == nil {
		return ""
	}
	// If already a string, keep it.
	if s, ok := v.(string); ok {
		return s
	}
	// Else JSON encode for readable text content.
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(b)
}