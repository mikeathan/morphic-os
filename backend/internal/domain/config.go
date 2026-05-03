package domain

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ProviderConfig holds configuration for a specific LLM provider.
type ProviderConfig struct {
	APIKey  string `json:"api_key" yaml:"api_key"`
	BaseURL string `json:"base_url,omitempty" yaml:"base_url,omitempty"`
	Model   string `json:"model,omitempty" yaml:"model,omitempty"`
}

// Config represents the global application configuration.
type Config struct {
	Port     string                    `json:"port" yaml:"port"`
	DBPath   string                    `json:"db_path" yaml:"db_path"`
	LLM      map[string]ProviderConfig `json:"llm" yaml:"llm"`
	Active   string                    `json:"active_llm" yaml:"active_llm"` // The provider to use by default
	Fallback string                    `json:"fallback_llm,omitempty" yaml:"fallback_llm,omitempty"` // The provider to use if Active fails
}

// LoadConfig loads the configuration from a JSON file.
func LoadConfig(path string) (*Config, error) {
	// Defaults
	cfg := &Config{
		Port:     "8080",
		DBPath:   "morphic-os.db",
		LLM:      make(map[string]ProviderConfig),
		Active:   "mock",
		Fallback: "",
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// If config doesn't exist, try to use ENV vars as fallback to maintain backwards compatibility during transition
			cfg.loadFromEnv()
			return cfg, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// We'll primarily support JSON for now since we don't have a YAML dependency installed yet.
	if filepath.Ext(path) == ".json" {
		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse json config: %w", err)
		}
	} else {
		return nil, fmt.Errorf("unsupported config file extension, use .json")
	}

	// Override with ENV if set (Env takes precedence for 12-factor compliance)
	cfg.loadFromEnv()

	return cfg, nil
}

func (c *Config) loadFromEnv() {
	if p := os.Getenv("PORT"); p != "" {
		c.Port = p
	}
	if d := os.Getenv("DB_PATH"); d != "" {
		c.DBPath = d
	}
	if a := os.Getenv("OPENAI_API_KEY"); a != "" {
		c.LLM["openai"] = ProviderConfig{
			APIKey:  a,
			BaseURL: os.Getenv("OPENAI_API_BASE"),
			Model:   os.Getenv("OPENAI_MODEL"),
		}
		c.Active = "openai"
	}
	if f := os.Getenv("FALLBACK_LLM"); f != "" {
		c.Fallback = f
	}
}
