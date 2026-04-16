package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunInit_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "portwatch.json")

	if err := runInit(out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(out); err != nil {
		t.Fatalf("expected config file to exist: %v", err)
	}
}

func TestRunInit_RefusesOverwrite(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "portwatch.json")

	// Create the file first
	if err := os.WriteFile(out, []byte(`{}`), 0644); err != nil {
		t.Fatal(err)
	}

	err := runInit(out)
	if err == nil {
		t.Fatal("expected error when file already exists")
	}
}

func TestRunInit_DefaultContent(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "portwatch.json")

	if err := runInit(out); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}

	content := string(data)
	if len(content) == 0 {
		t.Error("expected non-empty config file")
	}
	// Ensure localhost is present as default host
	if !contains(content, "localhost") {
		t.Error("expected localhost in default config")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
