package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"portwatch/internal/history"
)

func writeStabilityHistory(t *testing.T, dir string) {
	t.Helper()
	now := time.Now()
	s := history.NewMemoryStore()
	for i := 0; i < 4; i++ {
		s.Append("host-a", history.Entry{
			Timestamp: now.Add(time.Duration(i) * time.Hour),
			Ports:     []int{80, 443},
		})
	}
	if err := history.Save(dir, s); err != nil {
		t.Fatal(err)
	}
}

func TestRunStability_TextOutput(t *testing.T) {
	dir := t.TempDir()
	writeStabilityHistory(t, dir)

	out := captureOutput(t, func() {
		if err := runStability(dir, 2, "text"); err != nil {
			t.Fatal(err)
		}
	})
	if !containsStr(out, "host-a") {
		t.Errorf("expected host-a in output, got:\n%s", out)
	}
	if !containsStr(out, "100.0%") {
		t.Errorf("expected 100.0%% stability, got:\n%s", out)
	}
}

func TestRunStability_JSONOutput(t *testing.T) {
	dir := t.TempDir()
	writeStabilityHistory(t, dir)

	out := captureOutput(t, func() {
		if err := runStability(dir, 2, "json"); err != nil {
			t.Fatal(err)
		}
	})

	var results []map[string]interface{}
	if err := json.Unmarshal([]byte(out), &results); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out)
	}
	if len(results) == 0 {
		t.Error("expected at least one result")
	}
}

func TestRunStability_MissingHistory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nonexistent")
	out := captureOutput(t, func() {
		if err := runStability(dir, 2, "text"); err != nil {
			t.Fatal(err)
		}
	})
	if !containsStr(out, "No history") {
		t.Errorf("expected no-history message, got: %s", out)
	}
}

func TestRunStability_JSONStructure(t *testing.T) {
	dir := t.TempDir()
	writeStabilityHistory(t, dir)

	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	_ = runStability(dir, 2, "json")
	w.Close()
	os.Stdout = old

	buf := make([]byte, 4096)
	n, _ := r.Read(buf)
	var results []map[string]interface{}
	if err := json.Unmarshal(buf[:n], &results); err != nil {
		t.Fatalf("JSON parse error: %v", err)
	}
	for _, res := range results {
		if _, ok := res["Stability"]; !ok {
			t.Errorf("missing Stability field in %v", res)
		}
	}
}
