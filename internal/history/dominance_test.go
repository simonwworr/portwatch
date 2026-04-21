package history

import (
	"testing"
	"time"
)

func buildDominanceStore() Store {
	now := time.Now()
	s := NewMemoryStore()
	// hostA: ports 80, 443, 8080  — superset of hostB and hostC
	s.Append(Entry{Host: "hostA", Ports: []int{80, 443, 8080}, ScannedAt: now.Add(-2 * time.Hour)})
	s.Append(Entry{Host: "hostA", Ports: []int{80, 443, 8080}, ScannedAt: now.Add(-1 * time.Hour)})
	// hostB: ports 80, 443  — subset of hostA
	s.Append(Entry{Host: "hostB", Ports: []int{80, 443}, ScannedAt: now.Add(-2 * time.Hour)})
	s.Append(Entry{Host: "hostB", Ports: []int{80, 443}, ScannedAt: now.Add(-1 * time.Hour)})
	// hostC: ports 80  — subset of both hostA and hostB
	s.Append(Entry{Host: "hostC", Ports: []int{80}, ScannedAt: now.Add(-2 * time.Hour)})
	s.Append(Entry{Host: "hostC", Ports: []int{80}, ScannedAt: now.Add(-1 * time.Hour)})
	return s
}

func TestDominance_Basic(t *testing.T) {
	s := buildDominanceStore()
	results := Dominance(s, 1)
	if len(results) == 0 {
		t.Fatal("expected results, got none")
	}
	if results[0].Host != "hostA" {
		t.Errorf("expected hostA to be most dominant, got %s", results[0].Host)
	}
}

func TestDominance_Score(t *testing.T) {
	s := buildDominanceStore()
	results := Dominance(s, 1)
	var hostA DominanceResult
	for _, r := range results {
		if r.Host == "hostA" {
			hostA = r
		}
	}
	// hostA dominates both hostB and hostC → score 1.0
	if hostA.Score != 1.0 {
		t.Errorf("expected hostA score 1.0, got %.2f", hostA.Score)
	}
}

func TestDominance_MinPortsFilters(t *testing.T) {
	s := buildDominanceStore()
	// minPorts=3 should exclude hostB (2 ports) and hostC (1 port)
	results := Dominance(s, 3)
	if len(results) != 1 {
		t.Errorf("expected 1 result with minPorts=3, got %d", len(results))
	}
	if results[0].Host != "hostA" {
		t.Errorf("expected hostA, got %s", results[0].Host)
	}
}

func TestDominance_EmptyStore(t *testing.T) {
	s := NewMemoryStore()
	results := Dominance(s, 1)
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}

func TestDominance_NoDominance(t *testing.T) {
	now := time.Now()
	s := NewMemoryStore()
	// Disjoint port sets — neither dominates the other
	s.Append(Entry{Host: "hostA", Ports: []int{80}, ScannedAt: now})
	s.Append(Entry{Host: "hostB", Ports: []int{443}, ScannedAt: now})
	results := Dominance(s, 1)
	for _, r := range results {
		if r.Score != 0.0 {
			t.Errorf("expected score 0 for %s, got %.2f", r.Host, r.Score)
		}
	}
}
