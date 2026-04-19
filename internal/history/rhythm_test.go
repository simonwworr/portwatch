package history

import (
	"testing"
	"time"
)

func buildRhythmStore(intervals []time.Duration) Store {
	base := time.Now().Add(-time.Hour * 24)
	store := NewMemoryStore()
	t := base
	for _, d := range intervals {
		store.Append(Entry{Host: "host1", Timestamp: t, Ports: []int{80}})
		t = t.Add(d)
	}
	return store
}

func TestRhythm_Regular(t *testing.T) {
	intervals := []time.Duration{
		60 * time.Second,
		61 * time.Second,
		59 * time.Second,
		60 * time.Second,
	}
	store := buildRhythmStore(intervals)
	results := Rhythm(store, 3)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Regular {
		t.Errorf("expected regular rhythm")
	}
}

func TestRhythm_Irregular(t *testing.T) {
	intervals := []time.Duration{
		10 * time.Second,
		300 * time.Second,
		5 * time.Second,
		200 * time.Second,
	}
	store := buildRhythmStore(intervals)
	results := Rhythm(store, 3)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Regular {
		t.Errorf("expected irregular rhythm")
	}
}

func TestRhythm_InsufficientScans(t *testing.T) {
	store := buildRhythmStore([]time.Duration{60 * time.Second})
	results := Rhythm(store, 5)
	if len(results) != 0 {
		t.Errorf("expected no results, got %d", len(results))
	}
}

func TestRhythm_EmptyStore(t *testing.T) {
	store := NewMemoryStore()
	results := Rhythm(store, 2)
	if len(results) != 0 {
		t.Errorf("expected no results")
	}
}

func TestRhythm_ScanCount(t *testing.T) {
	intervals := []time.Duration{60 * time.Second, 60 * time.Second, 60 * time.Second}
	store := buildRhythmStore(intervals)
	results := Rhythm(store, 2)
	if results[0].ScanCount != 4 {
		t.Errorf("expected 4 scans, got %d", results[0].ScanCount)
	}
}
