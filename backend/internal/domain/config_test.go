package domain_test

import (
	"morphic-os/backend/internal/domain"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "morphic-config-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "config.json")
	configJSON := `{
		"port": "9090",
		"db_path": "test.db",
		"active_llm": "gemini",
		"llm": {
			"gemini": {
				"api_key": "test-key",
				"model": "test-model"
			}
		}
	}`

	if err := os.WriteFile(configPath, []byte(configJSON), 0644); err != nil {
		t.Fatalf("failed to write test config file: %v", err)
	}

	cfg, err := domain.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.Port != "9090" {
		t.Errorf("expected port '9090', got %q", cfg.Port)
	}
	if cfg.DBPath != "test.db" {
		t.Errorf("expected db_path 'test.db', got %q", cfg.DBPath)
	}
	if cfg.Active != "gemini" {
		t.Errorf("expected active_llm 'gemini', got %q", cfg.Active)
	}

	geminiCfg, ok := cfg.LLM["gemini"]
	if !ok {
		t.Fatalf("expected gemini config to exist")
	}
	if geminiCfg.APIKey != "test-key" {
		t.Errorf("expected api_key 'test-key', got %q", geminiCfg.APIKey)
	}
}
