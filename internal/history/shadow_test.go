package history

import (
	"testing"
	"time"
)

func buildShadowStore() Store {
	now := time.Now()
	s := NewMemoryStore()
	// hostA: port 80 appears 5 times, port 9999 appears only once
	for i := 0; i < 5; i++ {
		s.Append(Entry{Host: "hostA", Time: now.Add(time.Duration(i) * time.Hour), Ports: []int{80, 443}})
	}
	s.Append(Entry{Host: "hostA", Time: now.Add(10 * time.Hour), Ports: []int{80, 9999}})
	// hostB: all ports appear frequently
	for i := 0; i < 6; i++ {
		s.Append(Entry{Host: "hostB", Time: now.Add(time.Duration(i) * time.Hour), Ports: []int{22, 80}})
	}
	return s
}

func TestShadow_FindsTransientPort(t *testing.T) {
	s := buildShadowStore()
	results := Shadow(s, 3, 1)
	if len(results) != 1 {
		t.Fatalf("expected 1 shadow port, got %d", len(results))
	}
	if results[0].Port != 9999 {
		t.Errorf("expected port 9999, got %d", results[0].Port)
	}
	if results[0].Host != "hostA" {
		t.Errorf("expected hostA, got %s", results[0].Host)
	}
	if results[0].Appearances != 1 {
		t.Errorf("expected 1 appearance, got %d", results[0].Appearances)
	}
}

func TestShadow_NoTransientPorts(t *testing.T) {
	s := buildShadowStore()
	results := Shadow(s, 3, 0)
	if len(results) != 0 {
		t.Errorf("expected no shadow ports, got %d", len(results))
	}
}

func TestShadow_InsufficientScans(t *testing.T) {
	s := NewMemoryStore()
	now := time.Now()
	s.Append(Entry{Host: "hostC", Time: now, Ports: []int{1234}})
	results := Shadow(s, 5, 2)
	if len(results) != 0 {
		t.Errorf("expected no results for host with insufficient scans, got %d", len(results))
	}
}

func TestShadow_EmptyStore(t *testing.T) {
	s := NewMemoryStore()
	results := Shadow(s, 3, 1)
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}

func TestShadow_SortedByHostAndPort(t *testing.T) {
	now := time.Now()
	s := NewMemoryStore()
	for i := 0; i < 5; i++ {
		s.Append(Entry{Host: "alpha", Time: now.Add(time.Duration(i) * time.Hour), Ports: []int{80}})
		s.Append(Entry{Host: "beta", Time: now.Add(time.Duration(i) * time.Hour), Ports: []int{80}})
	}
	s.Append(Entry{Host: "beta", Time: now.Add(10 * time.Hour), Ports: []int{7777}})
	s.Append(Entry{Host: "alpha", Time: now.Add(10 * time.Hour), Ports: []int{5555}})
	results := Shadow(s, 3, 1)
	if len(results) < 2 {
		t.Fatalf("expected at least 2 results, got %d", len(results))
	}
	if results[0].Host > results[1].Host {
		t.Errorf("results not sorted by host: %s > %s", results[0].Host, results[1].Host)
	}
}
