package llm_test

import (
	"context"
	"fmt"
	"morphic-os/backend/internal/domain"
	"morphic-os/backend/internal/infrastructure/llm"
	"testing"
)

func TestFallbackAgent_SuccessOnPrimary(t *testing.T) {
	primary := llm.NewMockAgent()
	primary.EvaluateTaskFunc = func(ctx context.Context, task string, tools []*domain.Tool) (string, error) {
		return "primary response", nil
	}

	fallback := llm.NewMockAgent()
	fallback.EvaluateTaskFunc = func(ctx context.Context, task string, tools []*domain.Tool) (string, error) {
		return "fallback response", nil
	}

	agent := llm.NewFallbackAgent(primary, fallback)

	resp, err := agent.EvaluateTask(context.Background(), "task", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp != "primary response" {
		t.Errorf("expected primary response, got %s", resp)
	}
}

func TestFallbackAgent_FallbackOnError(t *testing.T) {
	primary := llm.NewMockAgent()
	primary.EvaluateTaskFunc = func(ctx context.Context, task string, tools []*domain.Tool) (string, error) {
		return "", fmt.Errorf("primary failed")
	}

	fallback := llm.NewMockAgent()
	fallback.EvaluateTaskFunc = func(ctx context.Context, task string, tools []*domain.Tool) (string, error) {
		return "fallback response", nil
	}

	agent := llm.NewFallbackAgent(primary, fallback)

	resp, err := agent.EvaluateTask(context.Background(), "task", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp != "fallback response" {
		t.Errorf("expected fallback response, got %s", resp)
	}
}

func TestFallbackAgent_BothFail(t *testing.T) {
	primary := llm.NewMockAgent()
	primary.EvaluateTaskFunc = func(ctx context.Context, task string, tools []*domain.Tool) (string, error) {
		return "", fmt.Errorf("primary failed")
	}

	fallback := llm.NewMockAgent()
	fallback.EvaluateTaskFunc = func(ctx context.Context, task string, tools []*domain.Tool) (string, error) {
		return "", fmt.Errorf("fallback failed")
	}

	agent := llm.NewFallbackAgent(primary, fallback)

	_, err := agent.EvaluateTask(context.Background(), "task", nil)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != "fallback failed" {
		t.Errorf("expected fallback error message, got %v", err)
	}
}
