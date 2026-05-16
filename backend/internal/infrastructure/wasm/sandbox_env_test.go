package wasm_test

import (
	"context"
	"morphic-os/backend/internal/domain"
	"morphic-os/backend/internal/infrastructure/wasm"
	"morphic-os/backend/internal/usecase"
	"testing"
)

func TestWazeroSandboxManager_ExecuteWASM_EnvVars(t *testing.T) {
	ctx := context.Background()
	eventBus := usecase.NewMemoryEventBus()
	sm := wasm.NewWazeroSandboxManager(ctx, eventBus)
	defer sm.Close(ctx)

	// A simple Go program that prints the value of "MY_ENV_VAR"
	sourceCode := `
	package main
	import (
		"fmt"
		"os"
	)
	func main() {
		fmt.Print(os.Getenv("MY_ENV_VAR"))
	}
	`

	wasmBytes, err := sm.CompileToWASM(ctx, "go", sourceCode)
	if err != nil {
		t.Fatalf("failed to compile: %v", err)
	}

	config := domain.SandboxConfig{
		EnvVars: map[string]string{
			"MY_ENV_VAR": "HelloSandbox",
		},
	}

	res, err := sm.ExecuteWASM(ctx, config, wasmBytes)
	if err != nil {
		t.Fatalf("failed to execute: %v", err)
	}

	if res.Stdout != "HelloSandbox" {
		t.Errorf("expected stdout 'HelloSandbox', got %q", res.Stdout)
	}
}
