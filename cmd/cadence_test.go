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

func writeCadenceHistory(t *testing.T, dir string) {
	t.Helper()
	now := time.Date(2024, 3, 1, 12, 0, 0, 0, time.UTC)
	s := history.NewMemoryStore()
	for i := 0; i < 5; i++ {
		s.Append(history.Entry{
			Host:  "host1",
			Time:  now.Add(time.Duration(i) * 60 * time.Second),
			Ports: []int{80, 443},
		})
	}
	if err := history.Save(dir, s); err != nil {
		t.Fatalf("save: %v", err)
	}
}

func TestRunCadence_TextOutput(t *testing.T) {
	dir := t.TempDir()
	writeCadenceHistory(t, dir)

	buf := &bytes.Buffer{}
	oldOut := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runCadence(dir, 2, 30.0, "text", false)
	w.Close()
	os.Stdout = oldOut
	buf.ReadFrom(r)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !containsStr(out, "host1") {
		t.Errorf("expected host1 in output, got: %s", out)
	}
	if !containsStr(out, "SCANS") {
		t.Errorf("expected header in output, got: %s", out)
	}
}

func TestRunCadence_JSONOutput(t *testing.T) {
	dir := t.TempDir()
	writeCadenceHistory(t, dir)

	tmpOut := filepath.Join(t.TempDir(), "out.json")
	f, _ := os.Create(tmpOut)
	oldOut := os.Stdout
	os.Stdout = f

	err := runCadence(dir, 2, 30.0, "json", false)
	f.Close()
	os.Stdout = oldOut

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(tmpOut)
	var results []map[string]interface{}
	if err := json.Unmarshal(data, &results); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(results) == 0 {
		t.Error("expected at least one result")
	}
	if results[0]["Host"] != "host1" {
		t.Errorf("expected host1, got %v", results[0]["Host"])
	}
}

func TestRunCadence_MissingHistory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nonexistent")
	err := runCadence(dir, 2, 30.0, "text", false)
	if err != nil {
		t.Fatalf("expected no error for missing dir, got: %v", err)
	}
}
