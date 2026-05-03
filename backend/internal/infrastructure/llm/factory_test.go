package llm_test

import (
	"morphic-os/backend/internal/domain"
	"morphic-os/backend/internal/infrastructure/llm"
	"testing"
)

func TestNewAgentFactory(t *testing.T) {
	t.Run("OpenAI Agent Creation", func(t *testing.T) {
		cfg := &domain.Config{
			Active: "openai",
			LLM: map[string]domain.ProviderConfig{
				"openai": {APIKey: "test-key"},
			},
		}
		agent := llm.NewAgentFactory(cfg)
		if _, ok := agent.(*llm.CloudProviderAgent); !ok {
			t.Errorf("expected *CloudProviderAgent, got %T", agent)
		}
	})

	t.Run("Fallback to Mock on Missing Key", func(t *testing.T) {
		cfg := &domain.Config{
			Active: "openai",
			LLM: map[string]domain.ProviderConfig{
				"openai": {APIKey: ""},
			},
		}
		agent := llm.NewAgentFactory(cfg)
		if _, ok := agent.(*llm.MockAgent); !ok {
			t.Errorf("expected *MockAgent due to fallback, got %T", agent)
		}
	})

	t.Run("Primary and Fallback Creation", func(t *testing.T) {
		cfg := &domain.Config{
			Active:   "openai",
			Fallback: "gemini",
			LLM: map[string]domain.ProviderConfig{
				"openai": {APIKey: "test-key"},
				"gemini": {APIKey: "test-key-gemini"},
			},
		}
		agent := llm.NewAgentFactory(cfg)
		if _, ok := agent.(*llm.FallbackAgent); !ok {
			t.Errorf("expected *FallbackAgent, got %T", agent)
		}
	})

	t.Run("Fallback Only (Primary Fails)", func(t *testing.T) {
		cfg := &domain.Config{
			Active:   "openai",
			Fallback: "gemini",
			LLM: map[string]domain.ProviderConfig{
				"openai": {APIKey: ""}, // Fails
				"gemini": {APIKey: "test-key-gemini"},
			},
		}
		agent := llm.NewAgentFactory(cfg)
		if _, ok := agent.(*llm.CloudProviderAgent); !ok {
			t.Errorf("expected *CloudProviderAgent (fallback), got %T", agent)
		}
	})

	t.Run("Mock Provider Explicitly", func(t *testing.T) {
		cfg := &domain.Config{
			Active: "mock",
		}
		agent := llm.NewAgentFactory(cfg)
		if _, ok := agent.(*llm.MockAgent); !ok {
			t.Errorf("expected *MockAgent, got %T", agent)
		}
	})
}
