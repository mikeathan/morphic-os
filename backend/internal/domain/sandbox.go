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

// SandboxConfig holds configuration for the sandbox execution environment.
type SandboxConfig struct {
	EnvVars         map[string]string // Environment variables to inject
	WorkspaceFSDir  string            // Host directory to mount for file system access
	TimeoutSeconds  int               // Execution timeout in seconds
}

// SandboxManager defines the interface for safely executing dynamically generated code.
// Instead of taking source code, it now takes compiled WebAssembly (WASM) bytes.
type SandboxManager interface {
	// ExecuteWASM runs the provided WebAssembly binary bytes within an isolated environment.
	// It accepts execution arguments and returns the standard output, standard error, and any execution error.
	ExecuteWASM(ctx context.Context, config SandboxConfig, wasmBytes []byte, args ...string) (*ExecutionResult, error)

	// CompileToWASM compiles the provided source code of the specified language into WebAssembly bytes.
	CompileToWASM(ctx context.Context, language string, sourceCode string) ([]byte, error)
}
