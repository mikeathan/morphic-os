package llm

import (
	"context"
	"encoding/json"
	"morphic-os/backend/internal/domain"
	"morphic-os/backend/internal/usecase"
)

// MockAgent is a stub implementation of usecase.Agent.
type MockAgent struct {
	EvaluateTaskFunc func(ctx context.Context, task string, tools []*domain.Tool) (string, error)
	FixToolFunc      func(ctx context.Context, task string, code string, errorMessage string) (string, error)
}

// NewMockAgent creates a new MockAgent.
func NewMockAgent() *MockAgent {
	return &MockAgent{}
}

// EvaluateTask calls the mock function, or returns a default direct response.
func (m *MockAgent) EvaluateTask(ctx context.Context, task string, tools []*domain.Tool) (string, error) {
	if m.EvaluateTaskFunc != nil {
		return m.EvaluateTaskFunc(ctx, task, tools)
	}

	response := usecase.AgentResponse{
		Action:   "direct_response",
		Response: "This is a mock response from the LLM agent.",
	}
	bytes, _ := json.Marshal(response)
	return string(bytes), nil
}

// FixTool calls the mock function, or returns a default forged tool response.
func (m *MockAgent) FixTool(ctx context.Context, task string, code string, errorMessage string) (string, error) {
	if m.FixToolFunc != nil {
		return m.FixToolFunc(ctx, task, code, errorMessage)
	}

	response := usecase.AgentResponse{
		Action:     "sys_forge_tool",
		SourceCode: code + "\n// Fixed tool code\n",
	}
	bytes, _ := json.Marshal(response)
	return string(bytes), nil
}
