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

func writeReachabilityHistory(t *testing.T, dir string) {
	t.Helper()
	now := time.Now()
	s := history.NewMemoryStore()
	s.Append("alpha", history.Entry{Timestamp: now.Add(-2 * time.Hour), Ports: []int{80, 443}})
	s.Append("alpha", history.Entry{Timestamp: now.Add(-1 * time.Hour), Ports: []int{80}})
	s.Append("alpha", history.Entry{Timestamp: now, Ports: []int{80, 443}})
	s.Append("beta", history.Entry{Timestamp: now.Add(-2 * time.Hour), Ports: []int{22}})
	s.Append("beta", history.Entry{Timestamp: now.Add(-1 * time.Hour), Ports: []int{}})
	s.Append("beta", history.Entry{Timestamp: now, Ports: []int{22}})
	if err := history.Save(dir, s); err != nil {
		t.Fatalf("save: %v", err)
	}
}

func TestRunReachability_TextOutput(t *testing.T) {
	dir := t.TempDir()
	writeReachabilityHistory(t, dir)

	out := captureOutput(t, func() {
		if err := runReachability(dir, 2, "text"); err != nil {
			t.Fatal(err)
		}
	})

	if !strings.Contains(out, "alpha") {
		t.Errorf("expected alpha in output, got:\n%s", out)
	}
	if !strings.Contains(out, "beta") {
		t.Errorf("expected beta in output, got:\n%s", out)
	}
	if !strings.Contains(out, "UPTIME%") {
		t.Errorf("expected header in output, got:\n%s", out)
	}
}

func TestRunReachability_JSONOutput(t *testing.T) {
	dir := t.TempDir()
	writeReachabilityHistory(t, dir)

	out := captureOutput(t, func() {
		if err := runReachability(dir, 2, "json"); err != nil {
			t.Fatal(err)
		}
	})

	var results []map[string]interface{}
	if err := json.Unmarshal([]byte(out), &results); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestRunReachability_MissingHistory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "no-such-dir")
	out := captureOutput(t, func() {
		if err := runReachability(dir, 2, "text"); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(out, "No history") {
		t.Errorf("expected no-history message, got: %s", out)
	}
}

func TestRunReachability_JSONStructure(t *testing.T) {
	dir := t.TempDir()
	writeReachabilityHistory(t, dir)

	out := captureOutput(t, func() {
		if err := runReachability(dir, 2, "json"); err != nil {
			t.Fatal(err)
		}
	})

	var results []map[string]interface{}
	_ = json.Unmarshal([]byte(out), &results)
	for _, r := range results {
		if _, ok := r["uptime"]; !ok {
			t.Errorf("missing uptime field in result: %v", r)
		}
		if _, ok := r["avg_ports"]; !ok {
			t.Errorf("missing avg_ports field in result: %v", r)
		}
	}
	_ = os.Getenv("") // suppress unused import warning
}
