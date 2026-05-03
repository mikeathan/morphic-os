package llm

import (
	"fmt"
	"log"
	"morphic-os/backend/internal/domain"
	"morphic-os/backend/internal/usecase"
)

// NewAgentFactory initializes the correct agent based on the global configuration.
// It handles resolving the primary agent and the fallback agent, ensuring a single
// point of truth for agent instantiation. It guarantees a valid agent is returned (defaults to MockAgent on total failure).
func NewAgentFactory(cfg *domain.Config) usecase.Agent {
	primary, err := createAgent(cfg.Active, cfg.LLM[cfg.Active])
	if err != nil {
		log.Printf("[LLM Factory] Failed to initialize primary agent '%s': %v.", cfg.Active, err)
	}

	var fallback usecase.Agent
	if cfg.Fallback != "" {
		var fErr error
		fallback, fErr = createAgent(cfg.Fallback, cfg.LLM[cfg.Fallback])
		if fErr != nil {
			log.Printf("[LLM Factory] Failed to initialize fallback agent '%s': %v.", cfg.Fallback, fErr)
		}
	}

	if primary != nil && fallback != nil {
		return NewFallbackAgent(primary, fallback)
	} else if primary != nil {
		return primary
	} else if fallback != nil {
		return fallback
	}

	log.Println("[LLM Factory] No valid primary or fallback agent configured. Falling back to MockAgent.")
	return NewMockAgent()
}

// createAgent returns the appropriate Agent implementation based on the provider string.
func createAgent(provider string, cfg domain.ProviderConfig) (usecase.Agent, error) {
	if cfg.APIKey == "" && provider != "mock" && provider != "" {
		return nil, fmt.Errorf("API key is required for provider %s", provider)
	}

	switch provider {
	case "openai", "openrouter", "mule", "nvidia", "gemini":
		// OpenRouter, Mule, Nvidia, and Gemini (via compatibility layer) API structures are OpenAI compatible.
		// We use the CloudProviderAgent for all of them provided the correct base URL is passed in.
		return NewCloudProviderAgent(cfg.APIKey, cfg.BaseURL, cfg.Model), nil
	case "mock":
		return NewMockAgent(), nil
	case "":
		return nil, fmt.Errorf("provider is empty")
	default:
		return nil, fmt.Errorf("unknown LLM provider: %s", provider)
	}
}
