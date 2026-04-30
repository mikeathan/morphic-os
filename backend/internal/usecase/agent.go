package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"morphic-os/backend/internal/domain"
	"time"
)

// Agent defines the interface for interacting with an LLM.
type Agent interface {
	// EvaluateTask analyzes the given task and returns a planned action
	// (either tool calls, a direct response, or a sys_forge_tool call).
	EvaluateTask(ctx context.Context, task string, tools []*domain.Tool) (string, error)
	// FixTool asks the LLM to fix the tool code based on the compilation/execution error.
	FixTool(ctx context.Context, task string, code string, errorMessage string) (string, error)
}

// AgentResponse structures the expected output from EvaluateTask or FixTool.
type AgentResponse struct {
	Action      string   `json:"action"`                // "direct_response", "sys_forge_tool", "tool_call"
	Response    string   `json:"response,omitempty"`    // For direct_response
	ToolName    string   `json:"tool_name,omitempty"`   // For sys_forge_tool or tool_call
	Description string   `json:"description,omitempty"` // For sys_forge_tool
	SourceCode  string   `json:"source_code,omitempty"` // For sys_forge_tool
	JSONSchema  string   `json:"json_schema,omitempty"` // For sys_forge_tool
	Arguments   []string `json:"arguments,omitempty"`   // For tool_call
}

// MorphicLoop orchestrates the execution flow.
type MorphicLoop struct {
	toolRepo domain.ToolRepository
	agent    Agent
	sandbox  domain.SandboxManager
}

// NewMorphicLoop creates a new MorphicLoop instance.
func NewMorphicLoop(toolRepo domain.ToolRepository, agent Agent, sandbox domain.SandboxManager) *MorphicLoop {
	return &MorphicLoop{
		toolRepo: toolRepo,
		agent:    agent,
		sandbox:  sandbox,
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
	responseStr, err := l.agent.EvaluateTask(ctx, task, activeTools)
	if err != nil {
		return "", err
	}

	var response AgentResponse
	if err := json.Unmarshal([]byte(responseStr), &response); err != nil {
		return "", fmt.Errorf("failed to parse agent response: %w", err)
	}

	// 3. Execution based on action
	switch response.Action {
	case "direct_response":
		return response.Response, nil
	case "sys_forge_tool":
		return l.forgeTool(ctx, task, response, activeTools)
	case "tool_call":
		return l.executeTool(ctx, response)
	default:
		return "", fmt.Errorf("unknown action type: %s", response.Action)
	}
}

func (l *MorphicLoop) executeTool(ctx context.Context, response AgentResponse) (string, error) {
	tool, err := l.toolRepo.GetByName(ctx, response.ToolName)
	if err != nil {
		return "", fmt.Errorf("failed to get tool %s from repository: %w", response.ToolName, err)
	}

	compiledWasm, err := l.sandbox.CompileGoToWASM(ctx, tool.SourceCode)
	if err != nil {
		return "", fmt.Errorf("failed to compile tool %s: %w", tool.Name, err)
	}

	execResult, err := l.sandbox.ExecuteWASM(ctx, compiledWasm, response.Arguments...)
	if err != nil {
		return "", fmt.Errorf("tool execution failed: %w", err)
	}

	if execResult.Stderr != "" {
		// Log warning or pass back
		fmt.Printf("Tool stderr: %s\n", execResult.Stderr)
	}

	return execResult.Stdout, nil
}

func (l *MorphicLoop) forgeTool(ctx context.Context, task string, response AgentResponse, activeTools []*domain.Tool) (string, error) {
	sourceCode := response.SourceCode
	var compiledWasm []byte
	var err error
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		compiledWasm, err = l.sandbox.CompileGoToWASM(ctx, sourceCode)
		if err == nil {
			break // Compilation succeeded
		}

		// Compilation failed, ask LLM to self-correct
		errorMessage := err.Error()
		fixedResponseStr, fixErr := l.agent.FixTool(ctx, task, sourceCode, errorMessage)
		if fixErr != nil {
			return "", fmt.Errorf("failed to ask agent to fix tool: %w", fixErr)
		}

		var fixedResponse AgentResponse
		if unmarshalErr := json.Unmarshal([]byte(fixedResponseStr), &fixedResponse); unmarshalErr != nil {
			return "", fmt.Errorf("failed to parse fixed agent response: %w", unmarshalErr)
		}
		sourceCode = fixedResponse.SourceCode
	}

	if err != nil {
		return "", fmt.Errorf("failed to compile tool after %d retries: %w", maxRetries, err)
	}
	_ = compiledWasm // Keep compiler happy for now, won't need to save compiled file yet

	// Compilation succeeded, save the tool
	newTool := &domain.Tool{
		ID:          response.ToolName, // Simplified ID assignment, maybe better to use uuid
		Name:        response.ToolName,
		Description: response.Description,
		JSONSchema:  response.JSONSchema,
		Language:    "go",
		SourceCode:  sourceCode,
		Active:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := l.toolRepo.Create(ctx, newTool); err != nil {
		return "", fmt.Errorf("failed to register new tool: %w", err)
	}

	// We call Execute again to fetch the latest active tools and have the agent decide the next step
	return l.Execute(ctx, task)
}
