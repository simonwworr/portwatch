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

func writeMaturityHistory(t *testing.T, dir string) {
	t.Helper()
	now := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	s := history.NewMemoryStore()
	for i := 0; i < 6; i++ {
		s.Append(history.Entry{
			Host:  "hostA",
			Time:  now.AddDate(0, 0, -60+i*10),
			Ports: []int{80, 443},
		})
	}
	for i := 0; i < 4; i++ {
		s.Append(history.Entry{
			Host:  "hostB",
			Time:  now.AddDate(0, 0, -20+i*5),
			Ports: []int{9000 + i},
		})
	}
	if err := history.Save(dir, s); err != nil {
		t.Fatalf("save history: %v", err)
	}
}

func TestRunMaturity_TextOutput(t *testing.T) {
	dir := t.TempDir()
	writeMaturityHistory(t, dir)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	if err := runMaturity(dir, 2, "text"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	out := buf.String()

	if !containsStr(out, "hostA") {
		t.Errorf("expected hostA in output, got:\n%s", out)
	}
	if !containsStr(out, "SCORE") {
		t.Errorf("expected header in output, got:\n%s", out)
	}
}

func TestRunMaturity_JSONOutput(t *testing.T) {
	dir := t.TempDir()
	writeMaturityHistory(t, dir)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	if err := runMaturity(dir, 2, "json"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)

	var results []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &results); err != nil {
		t.Fatalf("invalid JSON: %v\noutput: %s", err, buf.String())
	}
	if len(results) == 0 {
		t.Error("expected at least one result in JSON output")
	}
}

func TestRunMaturity_MissingHistory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nonexistent")
	if err := runMaturity(dir, 2, "text"); err != nil {
		t.Errorf("expected no error for missing history, got %v", err)
	}
}
