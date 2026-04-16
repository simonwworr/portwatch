package history

import (
	"testing"
	"time"
)

func TestPrune_MaxAge(t *testing.T) {
	l := &Log{}
	old := Entry{Timestamp: time.Now().UTC().Add(-48 * time.Hour), Host: "h1", Ports: []int{80}}
	recent := Entry{Timestamp: time.Now().UTC().Add(-1 * time.Hour), Host: "h1", Ports: []int{443}}
	l.Entries = []Entry{old, recent}

	removed := Prune(l, PruneOptions{MaxAge: 24 * time.Hour})
	if removed != 1 {
		t.Errorf("expected 1 removed, got %d", removed)
	}
	if len(l.Entries) != 1 || l.Entries[0].Ports[0] != 443 {
		t.Errorf("expected only recent entry to remain")
	}
}

func TestPrune_MaxEntries(t *testing.T) {
	l := &Log{}
	for i := 0; i < 5; i++ {
		l.Append("h1", []int{80 + i})
	}

	removed := Prune(l, PruneOptions{MaxEntries: 3})
	if removed != 2 {
		t.Errorf("expected 2 removed, got %d", removed)
	}
	if len(l.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(l.Entries))
	}
}

func TestPrune_NoOp(t *testing.T) {
	l := &Log{}
	l.Append("h1", []int{22})
	removed := Prune(l, PruneOptions{})
	if removed != 0 {
		t.Errorf("expected 0 removed, got %d", removed)
	}
}

func TestPrune_MultipleHosts_MaxEntries(t *testing.T) {
	l := &Log{}
	for i := 0; i < 4; i++ {
		l.Append("h1", []int{i})
		l.Append("h2", []int{i})
	}
	Prune(l, PruneOptions{MaxEntries: 2})
	if len(l.Entries) != 4 {
		t.Errorf("expected 4 entries (2 per host), got %d", len(l.Entries))
	}
}
