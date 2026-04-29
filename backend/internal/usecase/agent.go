package usecase

import (
	"context"
	"morphic-os/backend/internal/domain"
)

// Agent defines the interface for interacting with an LLM.
type Agent interface {
	// EvaluateTask analyzes the given task and returns a planned action
	// (either tool calls, a direct response, or a sys_forge_tool call).
	EvaluateTask(ctx context.Context, task string, tools []*domain.Tool) (string, error)
}

// MorphicLoop orchestrates the execution flow.
type MorphicLoop struct {
	toolRepo domain.ToolRepository
	agent    Agent
	// sandbox SandboxManager (to be added in Phase 2)
}

// NewMorphicLoop creates a new MorphicLoop instance.
func NewMorphicLoop(toolRepo domain.ToolRepository, agent Agent) *MorphicLoop {
	return &MorphicLoop{
		toolRepo: toolRepo,
		agent:    agent,
	}
}

// Execute handles a single user task.
func (l *MorphicLoop) Execute(ctx context.Context, task string) (string, error) {
	// 1. Context Assembly: Query SQLite for active tools
	activeTools, err := l.toolRepo.ListActive(ctx)
	if err != nil {
		return "", err
	}

	// 2. Evaluation: Ask LLM what to do
	// The LLM prompt and schema formatting logic will reside in the Agent implementation.
	response, err := l.agent.EvaluateTask(ctx, task, activeTools)
	if err != nil {
		return "", err
	}

	// TODO: Phase 3 implementation - Handle `sys_forge_tool` capability gaps
	// and execute existing tools.

	return response, nil
}
