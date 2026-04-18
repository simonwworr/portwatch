package history

import (
	"testing"
	"time"
)

func buildDiffStore() *Store {
	now := time.Now()
	s := &Store{}
	s.Append("host1", Entry{Time: now.Add(-2 * time.Minute), Ports: []int{80, 443}})
	s.Append("host1", Entry{Time: now.Add(-1 * time.Minute), Ports: []int{80, 443, 8080}})
	s.Append("host1", Entry{Time: now, Ports: []int{80}})
	return s
}

func TestDiffHistory_Basic(t *testing.T) {
	s := buildDiffStore()
	diffs := DiffHistory(s, "host1")

	if len(diffs) != 2 {
		t.Fatalf("expected 2 diffs, got %d", len(diffs))
	}

	if len(diffs[0].Opened) != 1 || diffs[0].Opened[0] != 8080 {
		t.Errorf("expected 8080 opened, got %v", diffs[0].Opened)
	}
	if len(diffs[0].Closed) != 0 {
		t.Errorf("expected no closed ports, got %v", diffs[0].Closed)
	}

	if len(diffs[1].Closed) != 2 {
		t.Errorf("expected 443 and 8080 closed, got %v", diffs[1].Closed)
	}
}

func TestDiffHistory_NoChanges(t *testing.T) {
	s := &Store{}
	now := time.Now()
	s.Append("host2", Entry{Time: now.Add(-time.Minute), Ports: []int{22, 80}})
	s.Append("host2", Entry{Time: now, Ports: []int{22, 80}})

	diffs := DiffHistory(s, "host2")
	if len(diffs) != 0 {
		t.Errorf("expected no diffs, got %d", len(diffs))
	}
}

func TestDiffHistory_InsufficientEntries(t *testing.T) {
	s := &Store{}
	s.Append("host3", Entry{Time: time.Now(), Ports: []int{80}})

	diffs := DiffHistory(s, "host3")
	if diffs != nil {
		t.Errorf("expected nil diffs for single entry, got %v", diffs)
	}
}

func TestDiffHistory_MissingHost(t *testing.T) {
	s := &Store{}
	diffs := DiffHistory(s, "ghost")
	if len(diffs) != 0 {
		t.Errorf("expected empty diffs for missing host")
	}
}
