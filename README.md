# Guardrail

> A policy-driven execution layer for securely exposing tools to AI agents via the Model Context Protocol (MCP).

Guardrail is an infrastructure library for building **secure, governed MCP servers in Go**.
It provides a deterministic execution pipeline that sits between AI agents and your internal systems — enforcing policies, logging activity, and protecting sensitive resources.

---

## Why Guardrail?

The Model Context Protocol (MCP) allows AI agents to call tools exposed by servers.

However, MCP alone does not provide:

- Fine-grained access control
- Rate limiting
- Execution timeouts
- Audit logging
- Policy enforcement
- Deterministic execution flow

Guardrail solves this by introducing a **policy-aware execution engine** that wraps tool invocation.

Instead of exposing raw functions to AI agents, Guardrail exposes:

> Governed capabilities.

---

## Architecture Overview

Guardrail separates responsibilities clearly:

```
MCP Client
    ↓
Transport (stdio / HTTP)
    ↓
MCP Adapter
    ↓
Execution Engine
    ↓
Tool
```

### Core Concepts

- **Tool** – A capability that can be executed.
- **Registry** – Stores and manages available tools.
- **Invocation** – Internal representation of a tool call.
- **Executor** – Runs tools through a policy pipeline.
- **Adapter (MCP)** – Bridges MCP protocol to Guardrail core.

The core execution engine is protocol-agnostic and does not depend on MCP.

---

## Design Principles

- Behavior-first tool abstraction
- Strict separation of protocol and execution logic
- Deterministic execution pipeline
- No global state
- Policy-driven extensibility
- Infrastructure-grade naming and structure

---

## Example Tool

```go
type HelloTool struct{}

func (HelloTool) Name() string { return "hello" }

func (HelloTool) Spec() tool.Spec {
    return tool.Spec{
        Title:       "Hello Tool",
        Description: "Greets a user by name",
        InputSchema: map[string]any{
            "type": "object",
            "properties": map[string]any{
                "name": map[string]any{
                    "type": "string",
                },
            },
            "required": []any{"name"},
        },
    }
}

func (HelloTool) Execute(ctx context.Context, input any) (any, error) {
    m, _ := input.(map[string]any)
    name, _ := m["name"].(string)
    return "Hello " + name, nil
}
```

---

## Registering Tools

```go
reg := tool.NewRegistry()
_ = reg.Add(HelloTool{})

exec := execution.NewExecutor(reg)
```

---

## Wiring MCP Adapter

```go
adapter, _ := mcpadapter.New("guardrail", "v0.1.0", reg, exec)
_ = adapter.RegisterAllTools()

adapter.Server().Run(context.Background(), &mcp.StdioTransport{})
```

Guardrail automatically:

- Exposes tools via `tools/list`
- Converts `ToolCall` → `Invocation`
- Executes tool via pipeline
- Returns `ToolResult`

---

## Execution Pipeline

When a tool is called:

1. ToolCall is converted to Invocation
2. Executor looks up tool in registry
3. Policies (future) run before execution
4. Tool executes
5. Policies run after execution
6. Result is returned to adapter

The execution engine is deterministic and safe to extend.

---

## Roadmap

- Policy engine (Timeout, Role, RateLimit)
- Audit logging
- OpenTelemetry integration
- Resource and Prompt support
- HTTP transport
- Structured input validation
- Tool annotations and metadata

---

## Status

Early development (MVP phase).

Core abstractions are in place:

- Tool
- Spec
- Registry
- Invocation
- Executor
- Minimal MCP adapter

---

## Philosophy

Guardrail treats AI tool exposure as infrastructure.

Instead of:

> “Here’s a function the AI can call.”

Guardrail provides:

> “Here’s a governed capability with enforced boundaries.”

