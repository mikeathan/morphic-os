package llm

import (
	"fmt"
	"morphic-os/backend/internal/domain"
	"morphic-os/backend/internal/usecase"
)

// NewAgentFactory returns the appropriate Agent implementation based on the provider string.
func NewAgentFactory(provider string, cfg domain.ProviderConfig) (usecase.Agent, error) {
	if cfg.APIKey == "" && provider != "mock" {
		return nil, fmt.Errorf("API key is required for provider %s", provider)
	}

	switch provider {
	case "openai", "openrouter", "mule", "nvidia", "gemini":
		// OpenRouter, Mule, Nvidia, and Gemini (via compatibility layer) API structures are OpenAI compatible.
		// We use the CloudProviderAgent for all of them provided the correct base URL is passed in.
		// For Gemini specifically, it requires the base URL pointing to the OpenAI compat endpoint, e.g.,
		// https://generativelanguage.googleapis.com/v1beta/openai/chat/completions
		return NewCloudProviderAgent(cfg.APIKey, cfg.BaseURL, cfg.Model), nil
	case "mock":
		return NewMockAgent(), nil
	default:
		return nil, fmt.Errorf("unknown LLM provider: %s", provider)
	}
}
