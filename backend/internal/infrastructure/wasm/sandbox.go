package wasm

import (
	"bytes"
	"context"
	"fmt"
	"morphic-os/backend/internal/domain"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"github.com/tetratelabs/wazero/sys"
)

// WazeroSandboxManager implements the domain.SandboxManager interface using wazero.
type WazeroSandboxManager struct {
	runtime wazero.Runtime
}

// NewWazeroSandboxManager creates a new SandboxManager powered by Wazero.
func NewWazeroSandboxManager(ctx context.Context) *WazeroSandboxManager {
	// Apply resource limits: Max memory set to something safe.
	// Note: Go programs compiled to wasip1 require at least 16-32MB just for the runtime.
	// 512 pages = 32MB. Let's use 1024 pages (64MB) to be safe for our sandbox.
	config := wazero.NewRuntimeConfig().WithMemoryLimitPages(1024)
	r := wazero.NewRuntimeWithConfig(ctx, config)

	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	return &WazeroSandboxManager{
		runtime: r,
	}
}

// CompileToWASM compiles the provided source code to WebAssembly based on the given language.
func (sm *WazeroSandboxManager) CompileToWASM(ctx context.Context, language string, sourceCode string) ([]byte, error) {
	switch language {
	case "go":
		return sm.compileGoToWASM(ctx, sourceCode)
	default:
		return nil, fmt.Errorf("unsupported language for WASM compilation: %s", language)
	}
}

// compileGoToWASM compiles the provided Go source code to WebAssembly using the local Go toolchain.
func (sm *WazeroSandboxManager) compileGoToWASM(ctx context.Context, sourceCode string) ([]byte, error) {
	// Create a temporary directory for the build
	tempDir, err := os.MkdirTemp("", "morphic-build-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir) // clean up

	// Write the source code to a file
	mainFile := filepath.Join(tempDir, "main.go")
	if err := os.WriteFile(mainFile, []byte(sourceCode), 0644); err != nil {
		return nil, fmt.Errorf("failed to write source file: %w", err)
	}

	// Initialize a Go module to support third-party dependencies
	// Naming it "mytool" as "tool" is reserved/causes issues in Go 1.21+
	modCmd := exec.CommandContext(ctx, "go", "mod", "init", "mytool")
	modCmd.Dir = tempDir
	if output, err := modCmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("failed to init module: %w\noutput: %s", err, string(output))
	}

	// Fetch dependencies
	// Set GOPATH to a valid temp dir, or keep the system default.
	// But it's important to set GOMODCACHE and other environments if isolated.
	// Here, we just use the default environment but set GOOS back to avoid tidy errors if any.
	tidyCmd := exec.CommandContext(ctx, "go", "mod", "tidy")
	tidyCmd.Dir = tempDir
	tidyEnv := append(os.Environ(), "GO111MODULE=on")
	tidyCmd.Env = tidyEnv
	if output, err := tidyCmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("failed to tidy module: %w\noutput: %s", err, string(output))
	}

	// Define the output path for the compiled WASM
	outputFile := filepath.Join(tempDir, "output.wasm")

	// Set up the go build command
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", outputFile, mainFile)
	buildCmd.Dir = tempDir
	buildCmd.Env = append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm")

	// Capture output for error reporting
	var stderrBuf bytes.Buffer
	buildCmd.Stderr = &stderrBuf

	if err := buildCmd.Run(); err != nil {
		return nil, fmt.Errorf("compilation failed: %w\ncompiler output: %s", err, stderrBuf.String())
	}

	// Read the compiled WASM binary
	wasmBytes, err := os.ReadFile(outputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read compiled wasm: %w", err)
	}

	return wasmBytes, nil
}

// Close cleans up the wazero runtime.
func (sm *WazeroSandboxManager) Close(ctx context.Context) error {
	return sm.runtime.Close(ctx)
}

// ExecuteWASM runs the WebAssembly code.
func (sm *WazeroSandboxManager) ExecuteWASM(ctx context.Context, sandboxConfig domain.SandboxConfig, wasmBytes []byte, args ...string) (*domain.ExecutionResult, error) {
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	cmdArgs := append([]string{"tool"}, args...)

	// Generate a unique name for the module to avoid collisions
	moduleName := uuid.New().String()

	config := wazero.NewModuleConfig().
		WithName(moduleName).
		WithStdout(&stdoutBuf).
		WithStderr(&stderrBuf).
		WithArgs(cmdArgs...)

	// Apply environment variables
	for k, v := range sandboxConfig.EnvVars {
		config = config.WithEnv(k, v)
	}

	// Mount the workspace filesystem if provided
	if sandboxConfig.VFSRepo != nil {
		vfs := NewWazeroVFS(ctx, sandboxConfig.VFSRepo, sandboxConfig.WorkspaceID)
		config = config.WithFSConfig(wazero.NewFSConfig().WithFSMount(vfs, "/"))
	} else if sandboxConfig.WorkspaceFSDir != "" {
		if err := os.MkdirAll(sandboxConfig.WorkspaceFSDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create workspace directory: %w", err)
		}
		config = config.WithFSConfig(wazero.NewFSConfig().WithDirMount(sandboxConfig.WorkspaceFSDir, "/"))
	}

	// Apply an execution timeout
	timeout := 30 * time.Second
	if sandboxConfig.TimeoutSeconds > 0 {
		timeout = time.Duration(sandboxConfig.TimeoutSeconds) * time.Second
	}
	execCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Compile the module first, so we don't compile concurrently if we cache later.
	compiled, err := sm.runtime.CompileModule(execCtx, wasmBytes)
	if err != nil {
		return &domain.ExecutionResult{
			Error: err,
		}, fmt.Errorf("failed to compile module: %w", err)
	}
	defer compiled.Close(execCtx)

	// Instantiate and execute the WASM module.
	mod, err := sm.runtime.InstantiateModule(execCtx, compiled, config)

	res := &domain.ExecutionResult{
		Stdout: stdoutBuf.String(),
		Stderr: stderrBuf.String(),
	}

	if mod != nil {
		defer mod.Close(execCtx)
	}

	if err != nil {
		// wazero runs `_start` on instantiation. It returns a sys.ExitError if the program calls os.Exit.
		// For Go wasip1, successful exits will be sys.ExitError with ExitCode 0.
		if exitErr, ok := err.(*sys.ExitError); ok {
			if exitErr.ExitCode() == 0 {
				// Normal, successful exit.
				res.Error = nil
				return res, nil
			}
		}

		// Some other error or non-zero exit code.
		res.Error = err
		return res, fmt.Errorf("wasm execution failed: %w", err)
	}

	return res, nil
}
