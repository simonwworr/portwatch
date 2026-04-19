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

func writeBurstHistory(t *testing.T, dir string) {
	t.Helper()
	s, err := history.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	s.Append("h1", history.Entry{Time: now.Add(-4 * time.Hour), Ports: []int{80}})
	s.Append("h1", history.Entry{Time: now.Add(-3 * time.Hour), Ports: []int{80, 443, 8080, 8443, 9090}})
	s.Append("h1", history.Entry{Time: now.Add(-2 * time.Hour), Ports: []int{80, 443, 8080, 8443, 9090, 3306}})
	if err := s.Save(dir); err != nil {
		t.Fatal(err)
	}
}

func TestRunBurst_TextOutput(t *testing.T) {
	dir := t.TempDir()
	writeBurstHistory(t, dir)

	out := captureOutput(t, func() {
		if err := runBurst("h1", dir, 2, 3, "text"); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(out, "h1") {
		t.Errorf("expected host in output, got: %s", out)
	}
}

func TestRunBurst_JSONOutput(t *testing.T) {
	dir := t.TempDir()
	writeBurstHistory(t, dir)

	out := captureOutput(t, func() {
		if err := runBurst("h1", dir, 2, 3, "json"); err != nil {
			t.Fatal(err)
		}
	})
	var results []map[string]interface{}
	if err := json.Unmarshal([]byte(out), &results); err != nil {
		t.Fatalf("invalid JSON: %v\noutput: %s", err, out)
	}
}

func TestRunBurst_MissingHistory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nonexistent")
	out := captureOutput(t, func() {
		_ = runBurst("h1", dir, 2, 3, "text")
	})
	_ = out // no crash expected
}

func TestRunBurst_NoBursts(t *testing.T) {
	dir := t.TempDir()
	s, _ := history.Load(dir)
	now := time.Now()
	s.Append("h1", history.Entry{Time: now.Add(-2 * time.Hour), Ports: []int{80}})
	s.Append("h1", history.Entry{Time: now.Add(-1 * time.Hour), Ports: []int{80, 443}})
	_ = s.Save(dir)

	out := captureOutput(t, func() {
		if err := runBurst("h1", dir, 2, 50, "text"); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(out, "No bursts") {
		t.Errorf("expected no-burst message, got: %s", out)
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
