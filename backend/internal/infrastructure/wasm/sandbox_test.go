package wasm_test

import (
	"context"
	"morphic-os/backend/internal/domain"
	"morphic-os/backend/internal/infrastructure/wasm"
	"os"
	"path/filepath"
	"testing"
)

func TestWazeroSandboxManager_ExecuteWASM(t *testing.T) {
	ctx := context.Background()

	// Load the pre-compiled WASM module
	wasmPath := filepath.Join("..", "..", "..", "testdata", "hello", "hello.wasm")
	wasmBytes, err := os.ReadFile(wasmPath)
	if err != nil {
		t.Fatalf("failed to read test wasm file: %v. Did you run 'GOOS=wasip1 GOARCH=wasm go build -o hello.wasm main.go' in testdata/hello?", err)
	}

	sm := wasm.NewWazeroSandboxManager(ctx)
	defer sm.Close(ctx)

	t.Run("Execute without args", func(t *testing.T) {
		res, err := sm.ExecuteWASM(ctx, domain.SandboxConfig{}, wasmBytes)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		expectedStdout := "Hello, World!\n"
		if res.Stdout != expectedStdout {
			t.Errorf("expected stdout %q, got %q", expectedStdout, res.Stdout)
		}
		if res.Stderr != "" {
			t.Errorf("expected empty stderr, got %q", res.Stderr)
		}
	})

	t.Run("Execute with args", func(t *testing.T) {
		res, err := sm.ExecuteWASM(ctx, domain.SandboxConfig{}, wasmBytes, "Morphic")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		expectedStdout := "Hello, Morphic!\n"
		if res.Stdout != expectedStdout {
			t.Errorf("expected stdout %q, got %q", expectedStdout, res.Stdout)
		}
	})

	t.Run("Execute multiple times to test collision handling", func(t *testing.T) {
		_, err1 := sm.ExecuteWASM(ctx, domain.SandboxConfig{}, wasmBytes, "First")
		if err1 != nil {
			t.Fatalf("expected no error on first execution, got %v", err1)
		}

		res2, err2 := sm.ExecuteWASM(ctx, domain.SandboxConfig{}, wasmBytes, "Second")
		if err2 != nil {
			t.Fatalf("expected no error on second execution, got %v", err2)
		}

		expectedStdout := "Hello, Second!\n"
		if res2.Stdout != expectedStdout {
			t.Errorf("expected stdout %q, got %q", expectedStdout, res2.Stdout)
		}
	})
}

func TestWazeroSandboxManager_CompileToWASM(t *testing.T) {
	ctx := context.Background()
	sm := wasm.NewWazeroSandboxManager(ctx)
	defer sm.Close(ctx)

	t.Run("Compile valid Go code", func(t *testing.T) {
		sourceCode := `package main
import "fmt"
func main() {
	fmt.Println("Dynamic compile test!")
}`

		wasmBytes, err := sm.CompileToWASM(ctx, "go", sourceCode)
		if err != nil {
			t.Fatalf("expected compilation to succeed, got error: %v", err)
		}

		if len(wasmBytes) == 0 {
			t.Fatalf("expected non-empty wasm bytes")
		}

		// Execute to verify it works
		res, err := sm.ExecuteWASM(ctx, domain.SandboxConfig{}, wasmBytes)
		if err != nil {
			t.Fatalf("failed to execute dynamically compiled wasm: %v", err)
		}

		expectedStdout := "Dynamic compile test!\n"
		if res.Stdout != expectedStdout {
			t.Errorf("expected stdout %q, got %q", expectedStdout, res.Stdout)
		}
	})

	t.Run("Compile invalid Go code", func(t *testing.T) {
		sourceCode := `package main
import "fmt"
func main() {
	fmt.Println("Missing bracket"
}`

		_, err := sm.CompileToWASM(ctx, "go", sourceCode)
		if err == nil {
			t.Fatalf("expected compilation to fail, but it succeeded")
		}
	})

	t.Run("Compile code with external dependency", func(t *testing.T) {
		sourceCode := `package main
import (
	"fmt"
	"github.com/google/uuid"
)
func main() {
	id := uuid.New()
	fmt.Printf("UUID length: %d\n", len(id.String()))
}`

		wasmBytes, err := sm.CompileToWASM(ctx, "go", sourceCode)
		if err != nil {
			t.Fatalf("expected compilation with dependencies to succeed, got error: %v", err)
		}

		res, err := sm.ExecuteWASM(ctx, domain.SandboxConfig{}, wasmBytes)
		if err != nil {
			t.Fatalf("failed to execute dynamically compiled wasm with dependency: %v", err)
		}

		expectedStdout := "UUID length: 36\n"
		if res.Stdout != expectedStdout {
			t.Errorf("expected stdout %q, got %q", expectedStdout, res.Stdout)
		}
	})
}
