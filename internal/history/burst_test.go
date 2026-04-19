package history

import (
	"testing"
	"time"
)

func buildBurstStore(t *testing.T) *Store {
	t.Helper()
	dir := t.TempDir()
	s, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	s.Append("host1", Entry{Time: now.Add(-4 * time.Hour), Ports: []int{80}})
	s.Append("host1", Entry{Time: now.Add(-3 * time.Hour), Ports: []int{80, 443}})
	s.Append("host1", Entry{Time: now.Add(-2 * time.Hour), Ports: []int{80, 443, 8080, 8443, 9090}})
	s.Append("host1", Entry{Time: now.Add(-1 * time.Hour), Ports: []int{80, 443, 8080, 8443, 9090, 3306}})
	return s
}

func TestDetectBursts_FindsBurst(t *testing.T) {
	s := buildBurstStore(t)
	results := DetectBursts(s, "host1", 2, 3)
	if len(results) == 0 {
		t.Fatal("expected burst results, got none")
	}
	if results[0].BurstSize < 3 {
		t.Errorf("expected burst size >= 3, got %d", results[0].BurstSize)
	}
}

func TestDetectBursts_NoBurst(t *testing.T) {
	dir := t.TempDir()
	s, _ := Load(dir)
	now := time.Now()
	s.Append("host1", Entry{Time: now.Add(-2 * time.Hour), Ports: []int{80}})
	s.Append("host1", Entry{Time: now.Add(-1 * time.Hour), Ports: []int{80, 443}})
	results := DetectBursts(s, "host1", 2, 10)
	if len(results) != 0 {
		t.Errorf("expected no burst, got %d", len(results))
	}
}

func TestDetectBursts_InsufficientEntries(t *testing.T) {
	dir := t.TempDir()
	s, _ := Load(dir)
	now := time.Now()
	s.Append("host1", Entry{Time: now, Ports: []int{80}})
	results := DetectBursts(s, "host1", 2, 1)
	if results != nil {
		t.Error("expected nil for insufficient entries")
	}
}

func TestDetectBursts_EmptyStore(t *testing.T) {
	dir := t.TempDir()
	s, _ := Load(dir)
	results := DetectBursts(s, "host1", 2, 1)
	if results != nil {
		t.Error("expected nil for empty store")
	}
}
