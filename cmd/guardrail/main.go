package main

import (
	"context"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	mcpadapter "github.com/jamalishaq/guardrail/adapter/mcp"
	"github.com/jamalishaq/guardrail/execution"
	"github.com/jamalishaq/guardrail/tool"
)

func main() {
	// STDIO SERVER WARNING:
	// Do NOT log to stdout; use stderr. (log package defaults to stderr.)
	log.SetOutput(os.Stderr)

	reg := tool.NewRegistry()
	// TODO: reg.Add(your tools...)

	exec := execution.NewExecutor(reg)

	adapter, err := mcpadapter.New("guardrail", "v0.1.0", reg, exec)
	if err != nil {
		log.Fatal(err)
	}
	if err := adapter.RegisterAllTools(); err != nil {
		log.Fatal(err)
	}

	// Run over stdin/stdout until disconnect.
	if err := adapter.Server().Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}