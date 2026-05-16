package usecase

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// CompilerService safely compiles code and captures logs for the LLM's self-correction loop.
type CompilerService interface {
	// Returns the compiled Wasm BLOB, the compiler's stdout/stderr logs, and any execution error.
	CompileToWasm(ctx context.Context, toolName, sourceCode, language string) (wasmBlob []byte, logs string, err error)
}

// LocalCompilerService implements CompilerService.
type LocalCompilerService struct{}

// NewLocalCompilerService creates a new LocalCompilerService.
func NewLocalCompilerService() *LocalCompilerService {
	return &LocalCompilerService{}
}

// CompileToWasm compiles the provided source code to WebAssembly.
func (s *LocalCompilerService) CompileToWasm(ctx context.Context, toolName, sourceCode, language string) ([]byte, string, error) {
	// For now, assume tinygo for compilation.
	if language != "tinygo" && language != "go" {
		// we default to assuming the source code is Go/TinyGo, but to be robust against "go" language strings:
		// the spec says "assume TinyGo for now".
	}

	// Create a temp directory: /tmp/morphic_builds/{toolName}
	// We use os.MkdirTemp to ensure we don't have collisions if multiple compiles happen for the same toolname.
	baseBuildDir := filepath.Join(os.TempDir(), "morphic_builds")
	if err := os.MkdirAll(baseBuildDir, 0755); err != nil {
		return nil, "", fmt.Errorf("failed to create base build dir: %w", err)
	}

	tempDir, err := os.MkdirTemp(baseBuildDir, toolName+"-*")
	if err != nil {
		return nil, "", fmt.Errorf("failed to create temp directory for tool %s: %w", toolName, err)
	}
	defer os.RemoveAll(tempDir) // cleanup

	// Write sourceCode to main.go
	mainFile := filepath.Join(tempDir, "main.go")
	if err := os.WriteFile(mainFile, []byte(sourceCode), 0644); err != nil {
		return nil, "", fmt.Errorf("failed to write main.go: %w", err)
	}

	// Output path for wasm
	outputFile := filepath.Join(tempDir, "out.wasm")

	// Setup os/exec.CommandContext using tinygo:
	// tinygo build -o /tmp/morphic_builds/{toolName}/out.wasm -target=wasi /tmp/morphic_builds/{toolName}/main.go
	cmd := exec.CommandContext(ctx, "tinygo", "build", "-o", outputFile, "-target=wasi", mainFile)

	// Capture cmd.Stdout and cmd.Stderr into a buffer.
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	// Run the compiler
	err = cmd.Run()

	// Combine stdout and stderr for logs
	logs := stdoutBuf.String()
	if stderrBuf.Len() > 0 {
		if len(logs) > 0 {
			logs += "\n"
		}
		logs += stderrBuf.String()
	}

	if err != nil {
		// If the compile fails, return the stderr string as logs so the Agentic Loop can read the syntax errors.
		return nil, logs, fmt.Errorf("compilation failed: %w", err)
	}

	// Read out.wasm into []byte
	wasmBlob, err := os.ReadFile(outputFile)
	if err != nil {
		return nil, logs, fmt.Errorf("failed to read compiled wasm blob: %w", err)
	}

	return wasmBlob, logs, nil
}
