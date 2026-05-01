package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"morphic-os/backend/internal/domain"
)

// GeminiAgent implements the usecase.Agent interface using the native Google Gemini API.
type GeminiAgent struct {
	apiKey string
	model  string
}

// NewGeminiAgent creates a new GeminiAgent.
func NewGeminiAgent(apiKey, model string) *GeminiAgent {
	if model == "" {
		model = "gemini-1.5-pro"
	}
	return &GeminiAgent{
		apiKey: apiKey,
		model:  model,
	}
}

func (a *GeminiAgent) EvaluateTask(ctx context.Context, task string, tools []*domain.Tool) (string, error) {
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
		systemPrompt += fmt.Sprintf("- Name: %s\n  Description: %s\n  JSONSchema: %s\n", t.Name, t.Description, t.JSONSchema)
	}

	systemPrompt += `
Respond ONLY with a valid JSON object matching this schema:
{
  "action": "direct_response" | "sys_forge_tool" | "tool_call",
  "response": "Text response if action is direct_response",
  "tool_name": "Name of the tool to forge or call",
  "description": "Description of the new tool (if forging)",
  "source_code": "The complete Go source code (if forging)",
  "json_schema": "JSON schema describing the tool's inputs (if forging)",
  "arguments": ["arg1", "arg2"] // Array of string arguments (if tool_call)
}`

	return a.callGemini(ctx, systemPrompt, task)
}

func (a *GeminiAgent) FixTool(ctx context.Context, task string, code string, errorMessage string) (string, error) {
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

	return a.callGemini(ctx, formattedPrompt, "Fix the code and output the JSON.")
}

func (a *GeminiAgent) callGemini(ctx context.Context, systemPrompt, userMessage string) (string, error) {
	baseURL := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", a.model, a.apiKey)

	reqBody := map[string]interface{}{
		"system_instruction": map[string]interface{}{
			"parts": []map[string]interface{}{
				{"text": systemPrompt},
			},
		},
		"contents": []map[string]interface{}{
			{
				"role": "user",
				"parts": []map[string]interface{}{
					{"text": userMessage},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"responseMimeType": "application/json",
			"temperature":      0.2,
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal gemini request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", baseURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create gemini request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("gemini api request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read gemini response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("gemini api returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var geminiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(bodyBytes, &geminiResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal gemini response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no valid content in gemini response")
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}
