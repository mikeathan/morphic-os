package llm_test

import (
	"morphic-os/backend/internal/domain"
	"morphic-os/backend/internal/infrastructure/llm"
	"testing"
)

func TestNewAgentFactory(t *testing.T) {
	t.Run("OpenAI Agent Creation", func(t *testing.T) {
		cfg := domain.ProviderConfig{
			APIKey: "test-key",
		}
		agent, err := llm.NewAgentFactory("openai", cfg)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if agent == nil {
			t.Fatalf("expected an agent, got nil")
		}
	})

	t.Run("Gemini Agent Creation", func(t *testing.T) {
		cfg := domain.ProviderConfig{
			APIKey: "test-key",
		}
		agent, err := llm.NewAgentFactory("gemini", cfg)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if agent == nil {
			t.Fatalf("expected an agent, got nil")
		}
	})

	t.Run("Missing API Key Error", func(t *testing.T) {
		cfg := domain.ProviderConfig{
			APIKey: "", // Missing
		}
		_, err := llm.NewAgentFactory("openai", cfg)
		if err == nil {
			t.Fatalf("expected error for missing API key")
		}
	})

	t.Run("Unknown Provider Error", func(t *testing.T) {
		cfg := domain.ProviderConfig{
			APIKey: "test-key",
		}
		_, err := llm.NewAgentFactory("unknown", cfg)
		if err == nil {
			t.Fatalf("expected error for unknown provider")
		}
	})

	t.Run("Mock Provider", func(t *testing.T) {
		cfg := domain.ProviderConfig{}
		agent, err := llm.NewAgentFactory("mock", cfg)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if agent == nil {
			t.Fatalf("expected an agent, got nil")
		}
	})
}
