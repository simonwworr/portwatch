package history

import (
	"testing"
	"time"
)

func buildResilienceStore() Store {
	now := time.Now()
	s := NewMemoryStore()

	// hostA: suffers a disruption and fully recovers.
	// scan 0: ports 80,443,8080
	// scan 1: ports 80 only  (big drop → disruption)
	// scan 2: ports 80,443,8080  (full recovery)
	s.Append("hostA", Entry{Timestamp: now.Add(-3 * time.Hour), Ports: []int{80, 443, 8080}})
	s.Append("hostA", Entry{Timestamp: now.Add(-2 * time.Hour), Ports: []int{80}})
	s.Append("hostA", Entry{Timestamp: now.Add(-1 * time.Hour), Ports: []int{80, 443, 8080}})

	// hostB: stable, no disruptions.
	s.Append("hostB", Entry{Timestamp: now.Add(-3 * time.Hour), Ports: []int{22, 80}})
	s.Append("hostB", Entry{Timestamp: now.Add(-2 * time.Hour), Ports: []int{22, 80}})
	s.Append("hostB", Entry{Timestamp: now.Add(-1 * time.Hour), Ports: []int{22, 80}})

	// hostC: disruption with no recovery.
	s.Append("hostC", Entry{Timestamp: now.Add(-3 * time.Hour), Ports: []int{80, 443}})
	s.Append("hostC", Entry{Timestamp: now.Add(-2 * time.Hour), Ports: []int{80}})
	// no scan 3 for hostC, so no recovery window

	return s
}

func TestResilience_Basic(t *testing.T) {
	s := buildResilienceStore()
	results := Resilience(s, 2, 0.4)
	if len(results) == 0 {
		t.Fatal("expected results, got none")
	}
}

func TestResilience_HostAFullRecovery(t *testing.T) {
	s := buildResilienceStore()
	results := Resilience(s, 2, 0.4)

	var r *ResilienceResult
	for i := range results {
		if results[i].Host == "hostA" {
			r = &results[i]
			break
		}
	}
	if r == nil {
		t.Fatal("hostA not found in results")
	}
	if r.Disruptions != 1 {
		t.Errorf("expected 1 disruption, got %d", r.Disruptions)
	}
	if r.AvgRecoveryRate != 1.0 {
		t.Errorf("expected avg recovery rate 1.0, got %f", r.AvgRecoveryRate)
	}
}

func TestResilience_HostBNoDisruptions(t *testing.T) {
	s := buildResilienceStore()
	results := Resilience(s, 2, 0.4)

	for _, r := range results {
		if r.Host == "hostB" {
			if r.Disruptions != 0 {
				t.Errorf("expected 0 disruptions for hostB, got %d", r.Disruptions)
			}
			return
		}
	}
	t.Error("hostB not found in results")
}

func TestResilience_MinScansFilters(t *testing.T) {
	s := buildResilienceStore()
	// hostC has only 2 entries; requiring 3 should exclude it.
	results := Resilience(s, 3, 0.4)

	for _, r := range results {
		if r.Host == "hostC" {
			t.Error("hostC should have been filtered out by minScans=3")
		}
	}
}

func TestResilience_EmptyStore(t *testing.T) {
	s := NewMemoryStore()
	results := Resilience(s, 2, 0.4)
	if len(results) != 0 {
		t.Errorf("expected no results for empty store, got %d", len(results))
	}
}
