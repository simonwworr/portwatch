package history

import (
	"testing"
	"time"
)

func buildCentroidStore() Store {
	now := time.Now()
	s := NewMemoryStore()
	// hostA: port 80 always, port 443 always, port 9000 once
	s.Append(Entry{Host: "hostA", Ports: []int{80, 443, 9000}, Time: now.Add(-4 * 24 * time.Hour)})
	s.Append(Entry{Host: "hostA", Ports: []int{80, 443}, Time: now.Add(-3 * 24 * time.Hour)})
	s.Append(Entry{Host: "hostA", Ports: []int{80, 443}, Time: now.Add(-2 * 24 * time.Hour)})
	s.Append(Entry{Host: "hostA", Ports: []int{80, 443}, Time: now.Add(-1 * 24 * time.Hour)})
	// hostB: single scan — should be excluded when minScans=2
	s.Append(Entry{Host: "hostB", Ports: []int{22, 80}, Time: now})
	return s
}

func TestCentroid_Basic(t *testing.T) {
	s := buildCentroidStore()
	results := Centroid(s, 2)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]
	if r.Host != "hostA" {
		t.Errorf("expected hostA, got %s", r.Host)
	}
}

func TestCentroid_MajorityPorts(t *testing.T) {
	s := buildCentroidStore()
	results := Centroid(s, 2)
	r := results[0]
	// port 80 and 443 appear in all 4 scans — should be in centroid
	// port 9000 appears in 1/4 — should NOT be in centroid
	portSet := map[int]bool{}
	for _, p := range r.Centroid {
		portSet[p] = true
	}
	if !portSet[80] {
		t.Error("expected port 80 in centroid")
	}
	if !portSet[443] {
		t.Error("expected port 443 in centroid")
	}
	if portSet[9000] {
		t.Error("port 9000 should not be in centroid")
	}
}

func TestCentroid_DeviationIsLow(t *testing.T) {
	s := buildCentroidStore()
	results := Centroid(s, 2)
	r := results[0]
	// Most scans match centroid closely — deviation should be low
	if r.Deviation > 0.5 {
		t.Errorf("expected low deviation, got %.3f", r.Deviation)
	}
}

func TestCentroid_MinScansFilters(t *testing.T) {
	s := buildCentroidStore()
	// minScans=5 should exclude hostA (4 scans)
	results := Centroid(s, 5)
	if len(results) != 0 {
		t.Errorf("expected 0 results with minScans=5, got %d", len(results))
	}
}

func TestCentroid_EmptyStore(t *testing.T) {
	s := NewMemoryStore()
	results := Centroid(s, 1)
	if len(results) != 0 {
		t.Errorf("expected 0 results for empty store, got %d", len(results))
	}
}
