package history

import (
	"testing"
	"time"
)

func buildSessionStore(t *testing.T) (string, []Session) {
	t.Helper()
	dir := t.TempDir()
	now := time.Now().UTC().Truncate(time.Second)
	sessions := []Session{
		{
			ID:        "s1",
			Host:      "host-a",
			StartedAt: now.Add(-2 * time.Hour),
			EndedAt:   now.Add(-1 * time.Hour),
			ScanCount: 12,
			PortsSeen: []int{80, 443},
		},
		{
			ID:        "s2",
			Host:      "host-b",
			StartedAt: now.Add(-30 * time.Minute),
			EndedAt:   now,
			ScanCount: 5,
			PortsSeen: []int{22},
		},
		{
			ID:        "s3",
			Host:      "host-a",
			StartedAt: now.Add(-10 * time.Minute),
			EndedAt:   now,
			ScanCount: 3,
			PortsSeen: []int{80},
		},
	}
	for _, s := range sessions {
		if err := AppendSession(dir, s); err != nil {
			t.Fatalf("AppendSession: %v", err)
		}
	}
	return dir, sessions
}

func TestSession_AppendAndLoad(t *testing.T) {
	dir, want := buildSessionStore(t)
	got, err := LoadSessions(dir)
	if err != nil {
		t.Fatalf("LoadSessions: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("expected %d sessions, got %d", len(want), len(got))
	}
}

func TestSession_LoadMissingFile(t *testing.T) {
	dir := t.TempDir()
	sessions, err := LoadSessions(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sessions) != 0 {
		t.Fatalf("expected empty slice, got %d", len(sessions))
	}
}

func TestSession_ForHost(t *testing.T) {
	dir, _ := buildSessionStore(t)
	all, _ := LoadSessions(dir)
	result := SessionsForHost(all, "host-a")
	if len(result) != 2 {
		t.Fatalf("expected 2 sessions for host-a, got %d", len(result))
	}
	if result[0].ID != "s1" || result[1].ID != "s3" {
		t.Errorf("unexpected order: %v", result)
	}
}

func TestSession_ForHost_NoEntries(t *testing.T) {
	dir, _ := buildSessionStore(t)
	all, _ := LoadSessions(dir)
	result := SessionsForHost(all, "host-z")
	if len(result) != 0 {
		t.Fatalf("expected 0, got %d", len(result))
	}
}

func TestSession_SaveLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	now := time.Now().UTC().Truncate(time.Second)
	s := Session{ID: "rt1", Host: "h", StartedAt: now, EndedAt: now.Add(time.Minute), ScanCount: 1, PortsSeen: []int{8080}}
	if err := AppendSession(dir, s); err != nil {
		t.Fatal(err)
	}
	got, err := LoadSessions(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ID != "rt1" || got[0].ScanCount != 1 {
		t.Errorf("round-trip mismatch: %+v", got)
	}
}
