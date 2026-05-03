package llm

import (
	"log"
	"morphic-os/backend/internal/domain"
	"morphic-os/backend/internal/usecase"
)

// NewAgentFactory initializes the correct agent based on the global configuration.
// It resolves the primary and fallback LLM configurations and returns a single
// configured Agent. It guarantees a valid agent is returned (defaults to MockAgent on total failure).
func NewAgentFactory(cfg *domain.Config) usecase.Agent {
	if cfg.Active == "mock" {
		return NewMockAgent()
	}

	var configs []domain.ProviderConfig

	if primaryCfg, exists := cfg.LLM[cfg.Active]; exists && primaryCfg.APIKey != "" {
		configs = append(configs, primaryCfg)
	} else if cfg.Active != "mock" {
		log.Printf("[LLM Factory] Primary agent config for '%s' is missing or lacks API key.", cfg.Active)
	}

	if cfg.Fallback != "" && cfg.Fallback != "mock" {
		if fallbackCfg, exists := cfg.LLM[cfg.Fallback]; exists && fallbackCfg.APIKey != "" {
			configs = append(configs, fallbackCfg)
		} else {
			log.Printf("[LLM Factory] Fallback agent config for '%s' is missing or lacks API key.", cfg.Fallback)
		}
	}

	if len(configs) > 0 {
		return NewCloudProviderAgent(configs...)
	}

	log.Println("[LLM Factory] No valid primary or fallback agent configured. Falling back to MockAgent.")
	return NewMockAgent()
}
