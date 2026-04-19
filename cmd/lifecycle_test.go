package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"portwatch/internal/history"
)

func writeLifecycleHistory(t *testing.T, dir string) {
	t.Helper()
	now := time.Now().UTC().Truncate(time.Second)
	s := history.NewMemoryStore()
	s.Append("h1", history.Entry{Host: "h1", Ports: []int{80, 443}, Time: now.Add(-3 * time.Hour)})
	s.Append("h1", history.Entry{Host: "h1", Ports: []int{80, 443, 22}, Time: now.Add(-2 * time.Hour)})
	s.Append("h1", history.Entry{Host: "h1", Ports: []int{80}, Time: now.Add(-1 * time.Hour)})
	if err := history.Save(dir, s); err != nil {
		t.Fatal(err)
	}
}

func TestRunLifecycle_TextOutput(t *testing.T) {
	dir := t.TempDir()
	writeLifecycleHistory(t, dir)

	out := captureOutput(t, func() {
		_ = runLifecycle(dir, "h1", "text")
	})
	if !strings.Contains(out, "PORT") {
		t.Errorf("expected header in output, got: %s", out)
	}
}

func TestRunLifecycle_JSONOutput(t *testing.T) {
	dir := t.TempDir()
	writeLifecycleHistory(t, dir)

	out := captureOutput(t, func() {
		_ = runLifecycle(dir, "h1", "json")
	})
	var result []map[string]interface{}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON: %v\noutput: %s", err, out)
	}
	if len(result) == 0 {
		t.Error("expected lifecycle events in JSON output")
	}
}

func TestRunLifecycle_MissingHistory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nonexistent")
	out := captureOutput(t, func() {
		_ = runLifecycle(dir, "h1", "text")
	})
	if !strings.Contains(out, "no history") {
		t.Errorf("expected 'no history' message, got: %s", out)
	}
}

func TestRunLifecycle_NoDataForHost(t *testing.T) {
	dir := t.TempDir()
	writeLifecycleHistory(t, dir)
	out := captureOutput(t, func() {
		_ = runLifecycle(dir, "unknown", "text")
	})
	if !strings.Contains(out, "no lifecycle") {
		t.Errorf("expected no lifecycle message, got: %s", out)
	}
}

func TestLifecycleFile_CreatedInDir(t *testing.T) {
	dir := t.TempDir()
	writeLifecycleHistory(t, dir)
	entries, _ := os.ReadDir(dir)
	if len(entries) == 0 {
		t.Error("expected history file in dir")
	}
	_ = filepath.Join(dir, entries[0].Name())
}
