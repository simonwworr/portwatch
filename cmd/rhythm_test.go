package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"portwatch/internal/history"
)

func writeRhythmHistory(t *testing.T, dir string) {
	t.Helper()
	store := history.NewMemoryStore()
	base := time.Now()
	for i := 0; i < 5; i++ {
		store.Append(history.Entry{
			Host:      "host1",
			Timestamp: base.Add(time.Duration(i) * 60 * time.Second),
			Ports:     []int{80, 443},
		})
	}
	if err := history.Save(dir, store); err != nil {
		t.Fatal(err)
	}
}

func TestRunRhythm_TextOutput(t *testing.T) {
	dir := t.TempDir()
	writeRhythmHistory(t, dir)
	out := captureOutput(t, func() {
		if err := runRhythm(dir, 3, "text"); err != nil {
			t.Fatal(err)
		}
	})
	if !containsStr(out, "host1") {
		t.Errorf("expected host1 in output, got: %s", out)
	}
	if !containsStr(out, "HOST") {
		t.Errorf("expected header in output")
	}
}

func TestRunRhythm_JSONOutput(t *testing.T) {
	dir := t.TempDir()
	writeRhythmHistory(t, dir)
	out := captureOutput(t, func() {
		if err := runRhythm(dir, 3, "json"); err != nil {
			t.Fatal(err)
		}
	})
	var results []map[string]interface{}
	if err := json.Unmarshal([]byte(out), &results); err != nil {
		t.Fatalf("invalid JSON: %v\noutput: %s", err, out)
	}
	if len(results) == 0 {
		t.Error("expected at least one result")
	}
}

func TestRunRhythm_MissingHistory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nonexistent")
	out := captureOutput(t, func() {
		if err := runRhythm(dir, 3, "text"); err != nil {
			t.Fatal(err)
		}
	})
	if !containsStr(out, "no history") {
		t.Errorf("expected no history message, got: %s", out)
	}
}

func TestRunRhythm_RegularFlag(t *testing.T) {
	dir := t.TempDir()
	writeRhythmHistory(t, dir)
	out := captureOutput(t, func() {
		_ = runRhythm(dir, 3, "text")
	})
	if !containsStr(out, "yes") && !containsStr(out, "no") {
		t.Errorf("expected regular column value in output: %s", out)
	}
}

func captureOutput(t *testing.T, fn func()) string {
	t.Helper()
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = old
	buf := make([]byte, 4096)
	n, _ := r.Read(buf)
	return string(buf[:n])
}
