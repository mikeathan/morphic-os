package usecase_test

import (
	"context"
	"encoding/json"
	"errors"
	"morphic-os/backend/internal/domain"
	"morphic-os/backend/internal/infrastructure/llm"
	"morphic-os/backend/internal/usecase"
	"testing"
)

// MockToolRepository implements domain.ToolRepository
type MockToolRepository struct {
	tools map[string]*domain.Tool
}

func NewMockToolRepository() *MockToolRepository {
	return &MockToolRepository{
		tools: make(map[string]*domain.Tool),
	}
}

func (m *MockToolRepository) Create(ctx context.Context, tool *domain.Tool) error {
	m.tools[tool.ID] = tool
	return nil
}

func (m *MockToolRepository) GetByID(ctx context.Context, id string) (*domain.Tool, error) {
	if t, ok := m.tools[id]; ok {
		return t, nil
	}
	return nil, errors.New("not found")
}

func (m *MockToolRepository) GetByName(ctx context.Context, name string) (*domain.Tool, error) {
	for _, t := range m.tools {
		if t.Name == name {
			return t, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *MockToolRepository) ListActive(ctx context.Context) ([]*domain.Tool, error) {
	var active []*domain.Tool
	for _, t := range m.tools {
		if t.Active {
			active = append(active, t)
		}
	}
	return active, nil
}

func (m *MockToolRepository) Update(ctx context.Context, tool *domain.Tool) error {
	m.tools[tool.ID] = tool
	return nil
}

func (m *MockToolRepository) Delete(ctx context.Context, id string) error {
	delete(m.tools, id)
	return nil
}

// MockSandboxManager implements domain.SandboxManager
type MockSandboxManager struct {
	compileError error
	execError    error
	execStdout   string
	execStderr   string
	compileCalls int
	execCalls    int
}

func (m *MockSandboxManager) CompileToWASM(ctx context.Context, language string, sourceCode string) ([]byte, error) {
	m.compileCalls++
	if m.compileError != nil {
		return nil, m.compileError
	}
	if language != "go" {
		return nil, errors.New("unsupported language in mock")
	}
	return []byte("fake-wasm-binary"), nil
}

func (m *MockSandboxManager) ExecuteWASM(ctx context.Context, wasmBytes []byte, args ...string) (*domain.ExecutionResult, error) {
	m.execCalls++
	if m.execError != nil {
		return nil, m.execError
	}
	return &domain.ExecutionResult{
		Stdout: m.execStdout,
		Stderr: m.execStderr,
	}, nil
}

func TestMorphicLoop_Execute_DirectResponse(t *testing.T) {
	repo := NewMockToolRepository()
	sandbox := &MockSandboxManager{}

	agent := llm.NewMockAgent()
	agent.EvaluateTaskFunc = func(ctx context.Context, task string, tools []*domain.Tool) (string, error) {
		resp := usecase.AgentResponse{
			Action:   "direct_response",
			Response: "Hello!",
		}
		b, _ := json.Marshal(resp)
		return string(b), nil
	}

	loop := usecase.NewMorphicLoop(repo, agent, sandbox)

	result, err := loop.Execute(context.Background(), "say hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Hello!" {
		t.Errorf("expected 'Hello!', got %q", result)
	}
}

func TestMorphicLoop_Execute_ForgeAndCall(t *testing.T) {
	repo := NewMockToolRepository()
	sandbox := &MockSandboxManager{
		execStdout: "Tool output",
	}

	callCount := 0
	agent := llm.NewMockAgent()
	agent.EvaluateTaskFunc = func(ctx context.Context, task string, tools []*domain.Tool) (string, error) {
		callCount++
		if callCount == 1 {
			// First call: Decide to forge a tool
			resp := usecase.AgentResponse{
				Action:      "sys_forge_tool",
				ToolName:    "my_tool",
				SourceCode:  "package main",
				Description: "test tool",
			}
			b, _ := json.Marshal(resp)
			return string(b), nil
		}

		// Second call: Call the forged tool
		resp := usecase.AgentResponse{
			Action:    "tool_call",
			ToolName:  "my_tool",
			Arguments: []string{"arg1"},
		}
		b, _ := json.Marshal(resp)
		return string(b), nil
	}

	loop := usecase.NewMorphicLoop(repo, agent, sandbox)

	result, err := loop.Execute(context.Background(), "do something")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Tool output" {
		t.Errorf("expected 'Tool output', got %q", result)
	}

	if sandbox.compileCalls != 1 {
		t.Errorf("expected 1 compile call, got %d", sandbox.compileCalls)
	}

	if sandbox.execCalls != 1 {
		t.Errorf("expected 1 exec call, got %d", sandbox.execCalls)
	}

	// Verify tool was saved to repo
	tool, err := repo.GetByName(context.Background(), "my_tool")
	if err != nil {
		t.Fatalf("tool not saved in repo")
	}
	if tool.ExecStats != 1 {
		t.Errorf("expected ExecStats to be 1, got %d", tool.ExecStats)
	}
}

func TestMorphicLoop_Execute_ExecutionSelfCorrection(t *testing.T) {
	repo := NewMockToolRepository()
	// Pre-seed a tool that will fail on execution
	repo.Create(context.Background(), &domain.Tool{
		ID:         "test-id",
		Name:       "failing_tool",
		Active:     true,
		SourceCode: "package main",
		WasmBinary: []byte("old-binary"),
	})

	sandbox := &MockSandboxManager{
		execStderr: "runtime error: index out of bounds", // First execution fails with stderr
	}

	agent := llm.NewMockAgent()
	evalCount := 0
	agent.EvaluateTaskFunc = func(ctx context.Context, task string, tools []*domain.Tool) (string, error) {
		evalCount++
		if evalCount == 1 {
			// First eval: Use the failing tool
			resp := usecase.AgentResponse{
				Action:   "tool_call",
				ToolName: "failing_tool",
			}
			b, _ := json.Marshal(resp)
			return string(b), nil
		}

		// Second eval (after correction and re-forge loop): Use the fixed tool
		resp := usecase.AgentResponse{
			Action:   "tool_call",
			ToolName: "failing_tool", // the name stays the same, but it's a new entity
		}
		b, _ := json.Marshal(resp)

		// Change sandbox to succeed for the next execution
		sandbox.execStderr = ""
		sandbox.execStdout = "fixed output"

		return string(b), nil
	}

	fixCalled := false
	agent.FixToolFunc = func(ctx context.Context, task string, code string, errorMessage string) (string, error) {
		fixCalled = true
		resp := usecase.AgentResponse{
			Action:     "sys_forge_tool",
			ToolName:   "failing_tool",
			SourceCode: "package main\n// fixed",
		}
		b, _ := json.Marshal(resp)
		return string(b), nil
	}

	loop := usecase.NewMorphicLoop(repo, agent, sandbox)

	result, err := loop.Execute(context.Background(), "do something")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "fixed output" {
		t.Errorf("expected 'fixed output', got %q", result)
	}

	if !fixCalled {
		t.Errorf("expected FixTool to be called")
	}

	if sandbox.compileCalls != 1 {
		t.Errorf("expected 1 compile call for the fixed tool, got %d", sandbox.compileCalls)
	}

	// Verify the original tool was deactivated and a new one was created
	activeCount := 0
	for _, tool := range repo.tools {
		if tool.Active {
			activeCount++
		}
	}
	if activeCount != 1 {
		t.Errorf("expected 1 active tool, got %d", activeCount)
	}
}
