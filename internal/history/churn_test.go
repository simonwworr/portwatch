package history

import (
	"testing"
	"time"
)

func buildChurnStore() Store {
	now := time.Now()
	s := NewMemoryStore()
	// host-a: 3 scans with varying ports
	s.Append("host-a", Entry{Timestamp: now.Add(-2 * time.Hour), Ports: []int{80, 443}})
	s.Append("host-a", Entry{Timestamp: now.Add(-1 * time.Hour), Ports: []int{80, 8080}})
	s.Append("host-a", Entry{Timestamp: now, Ports: []int{80, 443, 9090}})
	// host-b: 2 scans, no change
	s.Append("host-b", Entry{Timestamp: now.Add(-1 * time.Hour), Ports: []int{22, 80}})
	s.Append("host-b", Entry{Timestamp: now, Ports: []int{22, 80}})
	// host-c: only 1 scan (below minScans=2)
	s.Append("host-c", Entry{Timestamp: now, Ports: []int{443}})
	return s
}

func TestChurn_Basic(t *testing.T) {
	s := buildChurnStore()
	results := Churn(s, 2)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestChurn_HostAHasChurn(t *testing.T) {
	s := buildChurnStore()
	results := Churn(s, 2)
	// host-a should have the highest churn rate
	if results[0].Host != "host-a" {
		t.Errorf("expected host-a first, got %s", results[0].Host)
	}
	if results[0].Opened == 0 && results[0].Closed == 0 {
		t.Error("expected non-zero churn for host-a")
	}
}

func TestChurn_HostBNoChurn(t *testing.T) {
	s := buildChurnStore()
	results := Churn(s, 2)
	var hostB *ChurnResult
	for i := range results {
		if results[i].Host == "host-b" {
			hostB = &results[i]
		}
	}
	if hostB == nil {
		t.Fatal("host-b not found in results")
	}
	if hostB.Total != 0 {
		t.Errorf("expected 0 total churn for host-b, got %d", hostB.Total)
	}
	if hostB.ChurnRate != 0.0 {
		t.Errorf("expected 0.0 churn rate for host-b, got %f", hostB.ChurnRate)
	}
}

func TestChurn_BelowMinScansExcluded(t *testing.T) {
	s := buildChurnStore()
	results := Churn(s, 2)
	for _, r := range results {
		if r.Host == "host-c" {
			t.Error("host-c should be excluded (only 1 scan)")
		}
	}
}

func TestChurn_EmptyStore(t *testing.T) {
	s := NewMemoryStore()
	results := Churn(s, 2)
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}
