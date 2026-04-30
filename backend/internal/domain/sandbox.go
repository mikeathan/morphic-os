package domain

import (
	"context"
)

// ExecutionResult captures the output and status of a tool execution within the sandbox.
type ExecutionResult struct {
	Stdout string
	Stderr string
	Error  error
}

// SandboxManager defines the interface for safely executing dynamically generated code.
// Instead of taking source code, it now takes compiled WebAssembly (WASM) bytes.
type SandboxManager interface {
	// ExecuteWASM runs the provided WebAssembly binary bytes within an isolated environment.
	// It accepts execution arguments and returns the standard output, standard error, and any execution error.
	ExecuteWASM(ctx context.Context, wasmBytes []byte, args ...string) (*ExecutionResult, error)

	// CompileGoToWASM compiles the provided Go source code into WebAssembly bytes.
	CompileGoToWASM(ctx context.Context, sourceCode string) ([]byte, error)
}
