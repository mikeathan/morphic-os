package llm

import (
	"context"
	"log"
	"morphic-os/backend/internal/domain"
	"morphic-os/backend/internal/usecase"
)

// FallbackAgent implements the usecase.Agent interface and wraps a primary and fallback agent.
// It will attempt to use the primary agent, and if it fails, it will automatically switch to the fallback agent.
type FallbackAgent struct {
	Primary  usecase.Agent
	Fallback usecase.Agent
}

// NewFallbackAgent creates a new FallbackAgent.
func NewFallbackAgent(primary, fallback usecase.Agent) *FallbackAgent {
	return &FallbackAgent{
		Primary:  primary,
		Fallback: fallback,
	}
}

// EvaluateTask calls the primary agent, falling back to the fallback agent if an error occurs.
func (a *FallbackAgent) EvaluateTask(ctx context.Context, task string, tools []*domain.Tool) (string, error) {
	resp, err := a.Primary.EvaluateTask(ctx, task, tools)
	if err != nil && a.Fallback != nil {
		log.Printf("[FallbackAgent] Primary agent failed EvaluateTask: %v. Switching to Fallback.", err)
		return a.Fallback.EvaluateTask(ctx, task, tools)
	}
	return resp, err
}

// FixTool calls the primary agent, falling back to the fallback agent if an error occurs.
func (a *FallbackAgent) FixTool(ctx context.Context, task string, code string, errorMessage string) (string, error) {
	resp, err := a.Primary.FixTool(ctx, task, code, errorMessage)
	if err != nil && a.Fallback != nil {
		log.Printf("[FallbackAgent] Primary agent failed FixTool: %v. Switching to Fallback.", err)
		return a.Fallback.FixTool(ctx, task, code, errorMessage)
	}
	return resp, err
}
