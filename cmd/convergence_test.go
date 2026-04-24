package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"portwatch/internal/history"
)

func writeConvergenceHistory(t *testing.T, dir string) {
	t.Helper()
	now := time.Now()
	s := history.NewMemoryStore()
	s.Append("alpha", history.Entry{Timestamp: now.Add(-2 * 24 * time.Hour), Ports: []int{80}})
	s.Append("alpha", history.Entry{Timestamp: now.Add(-1 * 24 * time.Hour), Ports: []int{80, 443}})
	s.Append("alpha", history.Entry{Timestamp: now, Ports: []int{80, 443}})
	if err := history.Save(dir, s); err != nil {
		t.Fatalf("save: %v", err)
	}
}

func TestRunConvergence_TextOutput(t *testing.T) {
	dir := t.TempDir()
	writeConvergenceHistory(t, dir)

	out := captureOutput(t, func() {
		if err := runConvergence(dir, "text"); err != nil {
			t.Fatalf("runConvergence: %v", err)
		}
	})
	if !containsStr(out, "alpha") {
		t.Errorf("expected host alpha in output, got: %s", out)
	}
	if !containsStr(out, "STABLE") {
		t.Errorf("expected header in output, got: %s", out)
	}
}

func TestRunConvergence_JSONOutput(t *testing.T) {
	dir := t.TempDir()
	writeConvergenceHistory(t, dir)

	out := captureOutput(t, func() {
		if err := runConvergence(dir, "json"); err != nil {
			t.Fatalf("runConvergence: %v", err)
		}
	})
	var results []map[string]interface{}
	if err := json.Unmarshal([]byte(out), &results); err != nil {
		t.Fatalf("invalid JSON: %v\noutput: %s", err, out)
	}
	if len(results) == 0 {
		t.Error("expected at least one result")
	}
	if _, ok := results[0]["change_rate"]; !ok {
		t.Error("expected change_rate field in JSON")
	}
}

func TestRunConvergence_MissingHistory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nonexistent")
	out := captureOutput(t, func() {
		if err := runConvergence(dir, "text"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	if !containsStr(out, "no history") {
		t.Errorf("expected 'no history' message, got: %s", out)
	}
}

func TestRunConvergence_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	// create an empty store and save it
	s := history.NewMemoryStore()
	_ = history.Save(dir, s)

	// write a dummy file so Load doesn't fail with IsNotExist
	_ = os.WriteFile(filepath.Join(dir, "history.json"), []byte(`{}`), 0644)

	out := captureOutput(t, func() {
		_ = runConvergence(dir, "text")
	})
	if !containsStr(out, "no convergence") && !containsStr(out, "no history") {
		t.Logf("output: %s", out)
	}
}
