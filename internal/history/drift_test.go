package history

import (
	"testing"
	"time"
)

func buildDriftStore() Store {
	now := time.Now()
	s := NewMemoryStore()
	// host-a: baseline [80,443], latest [80,8080] => removed 443, added 8080
	s.Append("host-a", Entry{Ports: []int{80, 443}, ScannedAt: now.Add(-48 * time.Hour)})
	s.Append("host-a", Entry{Ports: []int{80, 8080}, ScannedAt: now})
	// host-b: no change
	s.Append("host-b", Entry{Ports: []int{22, 80}, ScannedAt: now.Add(-24 * time.Hour)})
	s.Append("host-b", Entry{Ports: []int{22, 80}, ScannedAt: now})
	// host-c: only one entry — should be skipped
	s.Append("host-c", Entry{Ports: []int{443}, ScannedAt: now})
	return s
}

func TestDrift_Basic(t *testing.T) {
	results := Drift(buildDriftStore())
	// host-c skipped, host-b score 0, host-a score > 0
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	// highest drift first
	if results[0].Host != "host-a" {
		t.Errorf("expected host-a first, got %s", results[0].Host)
	}
}

func TestDrift_AddedAndRemoved(t *testing.T) {
	results := Drift(buildDriftStore())
	var r DriftResult
	for _, x := range results {
		if x.Host == "host-a" {
			r = x
		}
	}
	if len(r.Added) != 1 || r.Added[0] != 8080 {
		t.Errorf("unexpected Added: %v", r.Added)
	}
	if len(r.Removed) != 1 || r.Removed[0] != 443 {
		t.Errorf("unexpected Removed: %v", r.Removed)
	}
	if len(r.Stable) != 1 || r.Stable[0] != 80 {
		t.Errorf("unexpected Stable: %v", r.Stable)
	}
}

func TestDrift_ScoreZeroWhenNoChange(t *testing.T) {
	results := Drift(buildDriftStore())
	for _, r := range results {
		if r.Host == "host-b" && r.DriftScore != 0 {
			t.Errorf("expected zero drift for host-b, got %f", r.DriftScore)
		}
	}
}

func TestDrift_SingleEntrySkipped(t *testing.T) {
	results := Drift(buildDriftStore())
	for _, r := range results {
		if r.Host == "host-c" {
			t.Error("host-c should have been skipped (only 1 entry)")
		}
	}
}

func TestDrift_EmptyStore(t *testing.T) {
	s := NewMemoryStore()
	results := Drift(s)
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}
