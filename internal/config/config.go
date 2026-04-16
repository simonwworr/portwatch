package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds the portwatch configuration.
type Config struct {
	Hosts      []string `json:"hosts"`
	Ports      []int    `json:"ports,omitempty"`
	Interval   int      `json:"interval_seconds,omitempty"`
	StateFile  string   `json:"state_file,omitempty"`
	WebhookURL string   `json:"webhook_url,omitempty"`
	LogFile    string   `json:"log_file,omitempty"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Interval:  60,
		StateFile: "portwatch_state.json",
	}
}

// Load reads a JSON config file from the given path.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	cfg := DefaultConfig()
	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return nil, fmt.Errorf("config: decode %q: %w", path, err)
	}
	if len(cfg.Hosts) == 0 {
		return nil, fmt.Errorf("config: at least one host is required")
	}
	return cfg, nil
}

// Save writes the config as JSON to the given path.
func Save(path string, cfg *Config) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("config: create %q: %w", path, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(cfg); err != nil {
		return fmt.Errorf("config: encode %q: %w", path, err)
	}
	return nil
}
