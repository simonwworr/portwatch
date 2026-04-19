package history

import (
	"testing"
	"time"
)

func buildSequenceStore() Store {
	now := time.Now()
	s := NewMemoryStore()
	// host-a: port 80 open in all 4 scans, port 443 only in 2
	s.Append("host-a", Entry{Time: now.Add(-4 * time.Hour), Ports: []int{80, 443}})
	s.Append("host-a", Entry{Time: now.Add(-3 * time.Hour), Ports: []int{80, 443}})
	s.Append("host-a", Entry{Time: now.Add(-2 * time.Hour), Ports: []int{80}})
	s.Append("host-a", Entry{Time: now.Add(-1 * time.Hour), Ports: []int{80}})
	// host-b: port 22 only once
	s.Append("host-b", Entry{Time: now.Add(-1 * time.Hour), Ports: []int{22}})
	return s
}

func TestSequence_FindsConsecutivePort(t *testing.T) {
	s := buildSequenceStore()
	res := Sequence(s, 3)
	if len(res) != 1 {
		t.Fatalf("expected 1 result, got %d", len(res))
	}
	if res[0].Host != "host-a" || res[0].Port != 80 {
		t.Errorf("unexpected entry: %+v", res[0])
	}
	if res[0].Runs != 4 {
		t.Errorf("expected 4 runs, got %d", res[0].Runs)
	}
}

func TestSequence_RespectsMinRuns(t *testing.T) {
	s := buildSequenceStore()
	res := Sequence(s, 2)
	// port 80 (4 runs) and port 443 (2 runs) should both appear
	ports := map[int]bool{}
	for _, r := range res {
		if r.Host == "host-a" {
			ports[r.Port] = true
		}
	}
	if !ports[80] || !ports[443] {
		t.Errorf("expected ports 80 and 443, got %v", ports)
	}
}

func TestSequence_InsufficientEntries(t *testing.T) {
	s := NewMemoryStore()
	s.Append("host-x", Entry{Time: time.Now(), Ports: []int{22}})
	res := Sequence(s, 2)
	if len(res) != 0 {
		t.Errorf("expected no results, got %d", len(res))
	}
}

func TestSequence_EmptyStore(t *testing.T) {
	s := NewMemoryStore()
	res := Sequence(s, 2)
	if len(res) != 0 {
		t.Errorf("expected empty result")
	}
}
