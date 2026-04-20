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

func writeSessionHistory(t *testing.T, dir string) {
	t.Helper()
	now := time.Now().UTC().Truncate(time.Second)
	sessions := []history.Session{
		{ID: "a1", Host: "192.168.1.1", StartedAt: now.Add(-1 * time.Hour), EndedAt: now, ScanCount: 6, PortsSeen: []int{80, 443}},
		{ID: "b1", Host: "10.0.0.1", StartedAt: now.Add(-30 * time.Minute), EndedAt: now, ScanCount: 3, PortsSeen: []int{22}},
	}
	for _, s := range sessions {
		if err := history.AppendSession(dir, s); err != nil {
			t.Fatalf("AppendSession: %v", err)
		}
	}
}

func TestRunSession_TextOutput(t *testing.T) {
	dir := t.TempDir()
	writeSessionHistory(t, dir)

	out := captureOutput(t, func() {
		if err := runSession(dir, "", "text"); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(out, "192.168.1.1") {
		t.Errorf("expected host in output, got: %s", out)
	}
	if !strings.Contains(out, "a1") {
		t.Errorf("expected session id in output, got: %s", out)
	}
}

func TestRunSession_JSONOutput(t *testing.T) {
	dir := t.TempDir()
	writeSessionHistory(t, dir)

	out := captureOutput(t, func() {
		if err := runSession(dir, "", "json"); err != nil {
			t.Fatal(err)
		}
	})
	var sessions []history.Session
	if err := json.Unmarshal([]byte(out), &sessions); err != nil {
		t.Fatalf("invalid JSON: %v\noutput: %s", err, out)
	}
	if len(sessions) != 2 {
		t.Errorf("expected 2 sessions, got %d", len(sessions))
	}
}

func TestRunSession_FilterByHost(t *testing.T) {
	dir := t.TempDir()
	writeSessionHistory(t, dir)

	out := captureOutput(t, func() {
		if err := runSession(dir, "10.0.0.1", "text"); err != nil {
			t.Fatal(err)
		}
	})
	if strings.Contains(out, "192.168.1.1") {
		t.Errorf("should not contain filtered-out host")
	}
	if !strings.Contains(out, "10.0.0.1") {
		t.Errorf("expected 10.0.0.1 in output")
	}
}

func TestRunSession_MissingHistory(t *testing.T) {
	dir := t.TempDir()
	out := captureOutput(t, func() {
		if err := runSession(dir, "", "text"); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(out, "no sessions") {
		t.Errorf("expected 'no sessions' message, got: %s", out)
	}
}

func TestSessionFile_CreatedInDir(t *testing.T) {
	dir := t.TempDir()
	writeSessionHistory(t, dir)
	path := filepath.Join(dir, "sessions.json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected sessions.json to exist at %s", path)
	}
}
