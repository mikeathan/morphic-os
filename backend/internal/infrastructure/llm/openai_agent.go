package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"morphic-os/backend/internal/domain"
)

// OpenAIAgent implements the usecase.Agent interface using the OpenAI API (or compatible local API).
type OpenAIAgent struct {
	apiKey  string
	baseURL string
	model   string
}

// NewOpenAIAgent creates a new OpenAIAgent.
func NewOpenAIAgent(apiKey, baseURL, model string) *OpenAIAgent {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1/chat/completions"
	} else if !strings.HasSuffix(baseURL, "/chat/completions") {
		// Ensure it points to the chat completions endpoint
		baseURL = strings.TrimRight(baseURL, "/") + "/chat/completions"
	}

	if model == "" {
		model = "gpt-4o"
	}
	return &OpenAIAgent{
		apiKey:  apiKey,
		baseURL: baseURL,
		model:   model,
	}
}

// EvaluateTask analyzes the given task and returns a planned action.
func (a *OpenAIAgent) EvaluateTask(ctx context.Context, task string, tools []*domain.Tool) (string, error) {
	systemPrompt := `You are Morphic-OS, an AI operating system.
Your goal is to complete the user's task.

You have access to a registry of WebAssembly tools.
If an active tool can complete the task, use "tool_call" action.
If NO tool can complete the task, you MUST dynamically create a new tool using the "sys_forge_tool" action.

When creating a new tool via "sys_forge_tool":
1. Write Go code (package main) that performs the task.
2. The code MUST read arguments from os.Args (os.Args[1] onwards).
3. The code MUST print the final result to stdout.
4. The code MUST compile cleanly with GOOS=wasip1 GOARCH=wasm.

Available Tools:
`
	for _, t := range tools {
		systemPrompt += fmt.Sprintf("- Name: %s\n  Description: %s\n", t.Name, t.Description)
	}

	systemPrompt += `
Respond ONLY with a valid JSON object matching this schema:
{
  "action": "direct_response" | "sys_forge_tool" | "tool_call",
  "response": "Text response if action is direct_response",
  "tool_name": "Name of the tool to forge or call",
  "description": "Description of the new tool (if forging)",
  "source_code": "The complete Go source code (if forging)",
  "arguments": ["arg1", "arg2"] // Array of string arguments (if tool_call)
}`

	return a.callOpenAI(ctx, systemPrompt, task)
}

func (a *OpenAIAgent) callOpenAI(ctx context.Context, systemPrompt, userMessage string) (string, error) {
	reqBody := map[string]interface{}{
		"model": a.model,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userMessage},
		},
		"response_format": map[string]string{"type": "json_object"},
		"temperature":     0.2,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal openai request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", a.baseURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create openai request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+a.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("openai api request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read openai response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("openai api returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var openaiResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(bodyBytes, &openaiResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal openai response: %w", err)
	}

	if len(openaiResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in openai response")
	}

	return openaiResp.Choices[0].Message.Content, nil
}

// FixTool asks the LLM to fix the tool code based on the compilation/execution error.
func (a *OpenAIAgent) FixTool(ctx context.Context, task string, code string, errorMessage string) (string, error) {
	systemPrompt := `You are Morphic-OS, an AI operating system.
Your goal is to fix the provided Go code that failed to compile or execute.

The code was meant to solve this task: "%s"

Here is the code that failed:
%s

Here is the error message:
%s

You MUST respond ONLY with a valid JSON object matching this schema:
{
  "action": "sys_forge_tool",
  "tool_name": "Name of the tool being fixed",
  "source_code": "The COMPLETE, FIXED Go source code (do not use markdown formatting, raw text only)"
}`

	formattedPrompt := fmt.Sprintf(systemPrompt, task, code, errorMessage)

	// We pass an empty string for userMessage here as the system prompt contains all the context
	return a.callOpenAI(ctx, formattedPrompt, "Fix the code and output the JSON.")
}
