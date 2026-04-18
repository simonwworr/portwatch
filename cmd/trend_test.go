package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"portwatch/internal/history"
)

func writeTrendHistory(t *testing.T, dir string) string {
	t.Helper()
	now := time.Now().UTC()
	s := &history.Store{}
	s.Append(history.Entry{Host: "10.0.0.1", ScannedAt: now.Add(-time.Hour), OpenPorts: []int{22, 80}})
	s.Append(history.Entry{Host: "10.0.0.1", ScannedAt: now, OpenPorts: []int{22}})
	p := filepath.Join(dir, "history.json")
	if err := history.Save(p, s); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunTrend_TextOutput(t *testing.T) {
	dir := t.TempDir()
	p := writeTrendHistory(t, dir)
	if err := runTrend("10.0.0.1", p, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunTrend_JSONOutput(t *testing.T) {
	dir := t.TempDir()
	p := writeTrendHistory(t, dir)

	// capture via temp file redirect trick — just ensure no error
	if err := runTrend("10.0.0.1", p, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunTrend_MissingHistory(t *testing.T) {
	if err := runTrend("10.0.0.1", "/nonexistent/path.json", false); err != nil {
		t.Fatalf("should handle missing file gracefully, got: %v", err)
	}
}

func TestRunTrend_JSONStructure(t *testing.T) {
	dir := t.TempDir()
	p := writeTrendHistory(t, dir)

	store, _ := history.Load(p)
	ht := history.Trend(store, "10.0.0.1")

	b, err := json.Marshal(ht)
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	if out["host"] != "10.0.0.1" {
		t.Errorf("expected host field, got %v", out["host"])
	}
	if _, ok := out["trends"]; !ok {
		t.Error("missing trends field")
	}
	_ = os.Remove(p)
}
