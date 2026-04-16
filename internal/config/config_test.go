package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_ValidConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	content := `{"hosts":["localhost"],"interval_seconds":30,"webhook_url":"http://example.com/hook"}`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Hosts) != 1 || cfg.Hosts[0] != "localhost" {
		t.Errorf("expected hosts [localhost], got %v", cfg.Hosts)
	}
	if cfg.Interval != 30 {
		t.Errorf("expected interval 30, got %d", cfg.Interval)
	}
	if cfg.WebhookURL != "http://example.com/hook" {
		t.Errorf("unexpected webhook url: %s", cfg.WebhookURL)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/config.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_NoHosts(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	os.WriteFile(path, []byte(`{"hosts":[]}`), 0644)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error when hosts is empty")
	}
}

func TestSaveLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	orig := &Config{
		Hosts:     []string{"host1", "host2"},
		Interval:  120,
		StateFile: "state.json",
	}
	if err := Save(path, orig); err != nil {
		t.Fatalf("save error: %v", err)
	}
	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if loaded.Interval != orig.Interval {
		t.Errorf("interval mismatch: got %d", loaded.Interval)
	}
	if len(loaded.Hosts) != 2 {
		t.Errorf("hosts mismatch: got %v", loaded.Hosts)
	}
}
