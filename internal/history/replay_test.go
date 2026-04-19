package history

import (
	"testing"
	"time"
)

func buildReplayStore(t *testing.T) Store {
	t.Helper()
	dir := t.TempDir()
	s, err := Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	now := time.Now()
	entries := []Entry{
		{Host: "h1", Ports: []int{80}, ScannedAt: now.Add(-3 * time.Hour)},
		{Host: "h1", Ports: []int{80, 443}, ScannedAt: now.Add(-2 * time.Hour)},
		{Host: "h1", Ports: []int{80, 443, 8080}, ScannedAt: now.Add(-1 * time.Hour)},
	}
	for _, e := range entries {
		if err := s.Append(e); err != nil {
			t.Fatalf("Append: %v", err)
		}
	}
	return s
}

func TestReplay_AllEntries(t *testing.T) {
	s := buildReplayStore(t)
	events, err := Replay(s, ReplayOptions{Host: "h1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(events))
	}
}

func TestReplay_Since(t *testing.T) {
	s := buildReplayStore(t)
	opts := ReplayOptions{Host: "h1", Since: time.Now().Add(-150 * time.Minute)}
	events, err := Replay(s, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
}

func TestReplay_Until(t *testing.T) {
	s := buildReplayStore(t)
	opts := ReplayOptions{Host: "h1", Until: time.Now().Add(-150 * time.Minute)}
	events, err := Replay(s, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
}

func TestReplay_NoEntries(t *testing.T) {
	s := buildReplayStore(t)
	events, err := Replay(s, ReplayOptions{Host: "unknown"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 0 {
		t.Fatalf("expected 0 events, got %d", len(events))
	}
}
