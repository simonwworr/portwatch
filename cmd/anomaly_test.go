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

func writeAnomalyHistory(t *testing.T, dir string) {
	t.Helper()
	s, _ := history.Load(dir)
	now := time.Now()
	for i := 0; i < 3; i++ {
		s.Append("host-a", history.Entry{Timestamp: now.Add(time.Duration(i) * time.Minute), Ports: []int{80}})
	}
	s.Append("host-a", history.Entry{Timestamp: now.Add(4 * time.Minute), Ports: []int{80, 9999}})
	if err := history.Save(dir, s); err != nil {
		t.Fatal(err)
	}
}

func TestRunAnomaly_TextOutput(t *testing.T) {
	dir := t.TempDir()
	writeAnomalyHistory(t, dir)

	out := captureOutput(t, func() {
		if err := runAnomaly(dir, 50.0, "text"); err != nil {
			t.Fatal(err)
		}
	})

	if !strings.Contains(out, "host-a") {
		t.Errorf("expected host-a in output, got: %s", out)
	}
	if !strings.Contains(out, "9999") {
		t.Errorf("expected port 9999 in output, got: %s", out)
	}
}

func TestRunAnomaly_JSONOutput(t *testing.T) {
	dir := t.TempDir()
	writeAnomalyHistory(t, dir)

	out := captureOutput(t, func() {
		if err := runAnomaly(dir, 50.0, "json"); err != nil {
			t.Fatal(err)
		}
	})

	var results []map[string]interface{}
	if err := json.Unmarshal([]byte(out), &results); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(results) == 0 {
		t.Error("expected at least one result")
	}
}

func TestRunAnomaly_MissingHistory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nonexistent")
	out := captureOutput(t, func() {
		if err := runAnomaly(dir, 50.0, "text"); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(out, "No history") {
		t.Errorf("expected no-history message, got: %s", out)
	}
}

func TestRunAnomaly_NoAnomalies(t *testing.T) {
	dir := t.TempDir()
	writeAnomalyHistory(t, dir)

	out := captureOutput(t, func() {
		_ = runAnomaly(dir, 5.0, "text")
	})
	if !strings.Contains(out, "No anomalies") {
		t.Errorf("expected no-anomalies message, got: %s", out)
	}
}

func captureOutput(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = old
	var buf strings.Builder
	buf.ReadFrom(r)
	return buf.String()
}
