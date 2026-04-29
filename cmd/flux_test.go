package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"portwatch/internal/history"
)

func writeFluxHistory(t *testing.T, dir string) {
	t.Helper()
	now := time.Now()
	st := history.NewMemoryStore()
	st.Append("alpha", history.Entry{Timestamp: now.Add(-4 * time.Hour), Ports: []int{80, 443}})
	st.Append("alpha", history.Entry{Timestamp: now.Add(-3 * time.Hour), Ports: []int{22, 8080}})
	st.Append("alpha", history.Entry{Timestamp: now.Add(-2 * time.Hour), Ports: []int{80, 443}})
	st.Append("alpha", history.Entry{Timestamp: now.Add(-1 * time.Hour), Ports: []int{22, 8080}})
	if err := history.Save(dir, st); err != nil {
		t.Fatalf("save history: %v", err)
	}
}

func TestRunFlux_TextOutput(t *testing.T) {
	dir := t.TempDir()
	writeFluxHistory(t, dir)

	out := captureOutput(t, func() {
		if err := runFlux(dir, 2, "text"); err != nil {
			t.Fatal(err)
		}
	})

	if !containsStr(out, "alpha") {
		t.Errorf("expected host alpha in output, got:\n%s", out)
	}
	if !containsStr(out, "FLUX") {
		t.Errorf("expected FLUX header in output, got:\n%s", out)
	}
}

func TestRunFlux_JSONOutput(t *testing.T) {
	dir := t.TempDir()
	writeFluxHistory(t, dir)

	out := captureOutput(t, func() {
		if err := runFlux(dir, 2, "json"); err != nil {
			t.Fatal(err)
		}
	})

	var results []map[string]interface{}
	if err := json.Unmarshal([]byte(out), &results); err != nil {
		t.Fatalf("invalid JSON: %v\noutput: %s", err, out)
	}
	if len(results) == 0 {
		t.Error("expected at least one result in JSON output")
	}
	if _, ok := results[0]["flux"]; !ok {
		t.Error("expected 'flux' key in JSON result")
	}
}

func TestRunFlux_MissingHistory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nonexistent")
	out := captureOutput(t, func() {
		if err := runFlux(dir, 2, "text"); err != nil {
			t.Fatal(err)
		}
	})
	if !containsStr(out, "No history") {
		t.Errorf("expected no-history message, got: %s", out)
	}
}

func TestRunFlux_JSONStructure(t *testing.T) {
	dir := t.TempDir()
	writeFluxHistory(t, dir)

	out := captureOutput(t, func() {
		_ = runFlux(dir, 2, "json")
	})

	var results []struct {
		Host     string  `json:"host"`
		Flux     float64 `json:"flux"`
		AvgDelta float64 `json:"avg_delta"`
		Scans    int     `json:"scans"`
	}
	if err := json.Unmarshal([]byte(out), &results); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if results[0].Host == "" {
		t.Error("expected non-empty host")
	}
	if results[0].Scans < 2 {
		t.Errorf("expected scans >= 2, got %d", results[0].Scans)
	}
	_ = os.Getenv("CI") // suppress unused import warning
}
