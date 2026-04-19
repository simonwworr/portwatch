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

func writePatternHistory(t *testing.T, dir string) {
	t.Helper()
	s := history.NewMemoryStore()
	now := time.Now()
	for i := 0; i < 4; i++ {
		s.Append("host-a", history.Entry{
			Timestamp: now.Add(time.Duration(i) * time.Hour),
			Ports:     []int{80, 443},
		})
	}
	if err := history.Save(dir, s); err != nil {
		t.Fatalf("save history: %v", err)
	}
}

func TestRunPattern_TextOutput(t *testing.T) {
	dir := t.TempDir()
	writePatternHistory(t, dir)

	out := captureOutput(t, func() {
		if err := runPattern(dir, 2, "text"); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(out, "host-a") {
		t.Errorf("expected host-a in output, got: %s", out)
	}
	if !strings.Contains(out, "80") {
		t.Errorf("expected port 80 in output, got: %s", out)
	}
}

func TestRunPattern_JSONOutput(t *testing.T) {
	dir := t.TempDir()
	writePatternHistory(t, dir)

	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	err := runPattern(dir, 2, "json")
	w.Close()
	os.Stdout = old
	if err != nil {
		t.Fatal(err)
	}
	var result []map[string]interface{}
	if err := json.NewDecoder(r).Decode(&result); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(result) == 0 {
		t.Error("expected json results")
	}
}

func TestRunPattern_MissingHistory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nonexistent")
	out := captureOutput(t, func() {
		if err := runPattern(dir, 1, "text"); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(out, "No pattern") {
		t.Errorf("expected no-data message, got: %s", out)
	}
}
