package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"portwatch/internal/history"
)

func writePressureHistory(t *testing.T, dir string) {
	t.Helper()
	now := time.Now()
	s := history.NewMemoryStore()
	for i := 0; i < 3; i++ {
		s.Append(history.Entry{Host: "alpha", Ports: []int{80, 443}, Timestamp: now.Add(time.Duration(i) * time.Hour)})
	}
	s.Append(history.Entry{Host: "beta", Ports: []int{22}, Timestamp: now})
	s.Append(history.Entry{Host: "beta", Ports: []int{22, 8080, 9090}, Timestamp: now.Add(time.Hour)})
	s.Append(history.Entry{Host: "beta", Ports: []int{22}, Timestamp: now.Add(2 * time.Hour)})
	if err := history.Save(dir, s); err != nil {
		t.Fatalf("save history: %v", err)
	}
}

func TestRunPressure_TextOutput(t *testing.T) {
	dir := t.TempDir()
	writePressureHistory(t, dir)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runPressure(dir, 2, "text")
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var buf bytes.Buffer
	buf.ReadFrom(r)
	out := buf.String()
	if !containsStr(out, "alpha") {
		t.Error("expected 'alpha' in output")
	}
	if !containsStr(out, "PRESSURE") {
		t.Error("expected header in output")
	}
}

func TestRunPressure_JSONOutput(t *testing.T) {
	dir := t.TempDir()
	writePressureHistory(t, dir)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runPressure(dir, 2, "json")
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var buf bytes.Buffer
	buf.ReadFrom(r)

	var results []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &results); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(results) == 0 {
		t.Error("expected non-empty JSON results")
	}
}

func TestRunPressure_MissingHistory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nonexistent")
	old := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w
	err := runPressure(dir, 2, "text")
	w.Close()
	os.Stdout = old
	if err != nil {
		t.Fatalf("expected nil error for missing dir, got %v", err)
	}
}

func TestRunPressure_JSONStructure(t *testing.T) {
	dir := t.TempDir()
	writePressureHistory(t, dir)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	_ = runPressure(dir, 2, "json")
	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	var results []map[string]interface{}
	_ = json.Unmarshal(buf.Bytes(), &results)
	for _, row := range results {
		if _, ok := row["Host"]; !ok {
			t.Error("missing Host field")
		}
		if _, ok := row["Pressure"]; !ok {
			t.Error("missing Pressure field")
		}
	}
}
