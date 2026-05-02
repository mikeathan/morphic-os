package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"morphic-os/backend/internal/domain"
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
	Language    string   `json:"language,omitempty"`    // For sys_forge_tool
	SourceCode  string   `json:"source_code,omitempty"` // For sys_forge_tool
	JSONSchema  string   `json:"json_schema,omitempty"` // For sys_forge_tool
	Arguments   []string `json:"arguments,omitempty"`   // For tool_call
}

// MorphicLoop orchestrates the execution flow.
type MorphicLoop struct {
	toolRepo      domain.ToolRepository
	workspaceRepo domain.WorkspaceRepository
	agent         Agent
	sandbox       domain.SandboxManager
	broadcastLog  func(msg string)
}

// NewMorphicLoop creates a new MorphicLoop instance.
func NewMorphicLoop(toolRepo domain.ToolRepository, workspaceRepo domain.WorkspaceRepository, agent Agent, sandbox domain.SandboxManager) *MorphicLoop {
	return &MorphicLoop{
		toolRepo:      toolRepo,
		workspaceRepo: workspaceRepo,
		agent:         agent,
		sandbox:       sandbox,
	}
}

// SetLogBroadcaster sets the function to call when a log event occurs.
func (l *MorphicLoop) SetLogBroadcaster(b func(msg string)) {
	l.broadcastLog = b
}

func (l *MorphicLoop) logEvent(level, message string) {
	if l.broadcastLog != nil {
		event := map[string]string{
			"type":    "log",
			"level":   level,
			"message": message,
		}
		bytes, err := json.Marshal(event)
		if err == nil {
			l.broadcastLog(string(bytes))
		}
	}
}

// Execute handles a single user task within a workspace.
func (l *MorphicLoop) Execute(ctx context.Context, workspaceID string, task string) (string, error) {
	l.logEvent("INFO", fmt.Sprintf("Received task in workspace %s: %q", workspaceID, task))

	// 1. Context Assembly: Query SQLite for active tools
	l.logEvent("EVAL", "Assembling context...")
	activeTools, err := l.toolRepo.ListActive(ctx, workspaceID)
	if err != nil {
		return "", err
	}

	// Implement Context Management: limit tools to prevent exceeding LLM context window.
	const maxTools = 10
	var contextTools []*domain.Tool
	if len(activeTools) > maxTools {
		contextTools = activeTools[:maxTools]
		l.logEvent("EVAL", fmt.Sprintf("Found %d active tools. Limiting context to %d tools.", len(activeTools), maxTools))
	} else {
		contextTools = activeTools
		l.logEvent("EVAL", fmt.Sprintf("Found %d active tools.", len(activeTools)))
	}

	// 2. Evaluation: Ask LLM what to do
	// The LLM prompt and schema formatting logic will reside in the Agent implementation.
	responseStr, err := l.agent.EvaluateTask(ctx, task, contextTools)
	if err != nil {
		return "", err
	}

	var response AgentResponse
	if err := json.Unmarshal([]byte(responseStr), &response); err != nil {
		return "", fmt.Errorf("failed to parse agent response: %w", err)
	}

	l.logEvent("EVAL", fmt.Sprintf("Agent decided action: %s", response.Action))

	// 3. Execution based on action
	switch response.Action {
	case "direct_response":
		return response.Response, nil
	case "sys_forge_tool":
		return l.forgeTool(ctx, workspaceID, task, response, activeTools)
	case "tool_call":
		return l.executeTool(ctx, workspaceID, task, response)
	default:
		return "", fmt.Errorf("unknown action type: %s", response.Action)
	}
}

func (l *MorphicLoop) executeTool(ctx context.Context, workspaceID string, task string, response AgentResponse) (string, error) {
	l.logEvent("EXEC", fmt.Sprintf("Executing tool %s", response.ToolName))
	tool, err := l.toolRepo.GetByName(ctx, workspaceID, response.ToolName)
	if err != nil {
		return "", fmt.Errorf("failed to get tool %s from repository: %w", response.ToolName, err)
	}

	compiledWasm := tool.WasmBinary
	if len(compiledWasm) == 0 {
		// Fallback: compile if we don't have the binary cached
		compiledWasm, err = l.sandbox.CompileToWASM(ctx, tool.Language, tool.SourceCode)
		if err != nil {
			return "", fmt.Errorf("failed to compile tool %s: %w", tool.Name, err)
		}
		// Save compiled binary back to db if needed
		tool.WasmBinary = compiledWasm
		_ = l.toolRepo.Update(ctx, tool)
	}

	sandboxConfig := domain.SandboxConfig{
		TimeoutSeconds: 30, // Default timeout
	}

	if l.workspaceRepo != nil && workspaceID != "" && workspaceID != "default" {
		workspace, err := l.workspaceRepo.GetByID(ctx, workspaceID)
		if err == nil {
			sandboxConfig.EnvVars = workspace.EnvVars
			baseDir := os.Getenv("WORKSPACE_BASE_DIR")
			if baseDir == "" {
				baseDir = "/tmp/morphic-os/workspaces/"
			}
			// Only use the base of the ID to avoid path traversal just in case
			sandboxConfig.WorkspaceFSDir = filepath.Join(baseDir, filepath.Base(workspace.ID))
		}
	}

	execResult, err := l.sandbox.ExecuteWASM(ctx, sandboxConfig, compiledWasm, response.Arguments...)

	// Check for execution errors or significant stderr indicating a failure
	if err != nil || execResult.Stderr != "" {
		errorMsg := ""
		if err != nil {
			errorMsg = err.Error()
		} else {
			errorMsg = fmt.Sprintf("tool executed but wrote to stderr: %s", execResult.Stderr)
		}

		l.logEvent("ERROR", fmt.Sprintf("Execution of tool %s failed: %v. Attempting self-correction.", tool.Name, errorMsg))

		// Attempt self-correction by asking the LLM to fix the execution issue
		fixedResponseStr, fixErr := l.agent.FixTool(ctx, task, tool.SourceCode, errorMsg)
		if fixErr != nil {
			return "", fmt.Errorf("tool execution failed and self-correction failed: %w", fixErr)
		}

		var fixedResponse AgentResponse
		if unmarshalErr := json.Unmarshal([]byte(fixedResponseStr), &fixedResponse); unmarshalErr != nil {
			return "", fmt.Errorf("failed to parse fixed agent response during execution correction: %w", unmarshalErr)
		}

		if fixedResponse.Description == "" {
			fixedResponse.Description = tool.Description
		}
		if fixedResponse.JSONSchema == "" {
			fixedResponse.JSONSchema = tool.JSONSchema
		}
		fixedResponse.ToolName = tool.Name

		tool.Active = false
		_ = l.toolRepo.Update(ctx, tool)

		return l.forgeTool(ctx, workspaceID, task, fixedResponse, nil)
	}

	// Increment ExecStats and save
	tool.ExecStats++
	tool.UpdatedAt = time.Now()
	if updateErr := l.toolRepo.Update(ctx, tool); updateErr != nil {
		l.logEvent("ERROR", fmt.Sprintf("Failed to update execution stats for tool %s: %v", tool.Name, updateErr))
	}

	return execResult.Stdout, nil
}

func (l *MorphicLoop) forgeTool(ctx context.Context, workspaceID string, task string, response AgentResponse, activeTools []*domain.Tool) (string, error) {
	l.logEvent("FORGE", fmt.Sprintf("Forging tool %s", response.ToolName))
	sourceCode := response.SourceCode
	language := response.Language
	if language == "" {
		language = "go"
	}
	var compiledWasm []byte
	var err error
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		l.logEvent("FORGE", fmt.Sprintf("Compiling tool to WebAssembly (Attempt %d)", i+1))
		compiledWasm, err = l.sandbox.CompileToWASM(ctx, language, sourceCode)
		if err == nil {
			l.logEvent("SUCCESS", "Compilation succeeded")
			break // Compilation succeeded
		}

		// Compilation failed, ask LLM to self-correct
		errorMessage := err.Error()
		l.logEvent("ERROR", fmt.Sprintf("Compilation failed: %v", errorMessage))
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

	// Compilation succeeded, save the tool
	l.logEvent("INFO", fmt.Sprintf("Saving tool %s to registry", response.ToolName))
	newTool := &domain.Tool{
		ID:          uuid.New().String(),
		WorkspaceID: workspaceID,
		Name:        response.ToolName,
		Description: response.Description,
		JSONSchema:  response.JSONSchema,
		Language:    language,
		SourceCode:  sourceCode,
		WasmBinary:  compiledWasm,
		Active:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := l.toolRepo.Create(ctx, newTool); err != nil {
		return "", fmt.Errorf("failed to register new tool: %w", err)
	}

	// We call Execute again to fetch the latest active tools and have the agent decide the next step
	return l.Execute(ctx, workspaceID, task)
}
