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
	case "openai", "openrouter", "mule", "nvidia":
		// OpenRouter, Mule, and Nvidia API structures are often OpenAI compatible.
		// For now, we will reuse the OpenAIAgent structure for all of them
		// provided the correct base URL is passed in.
		return NewOpenAIAgent(cfg.APIKey, cfg.BaseURL, cfg.Model), nil
	case "gemini":
		return NewGeminiAgent(cfg.APIKey, cfg.Model), nil
	case "mock":
		return NewMockAgent(), nil
	default:
		return nil, fmt.Errorf("unknown LLM provider: %s", provider)
	}
}
