package history

import (
	"testing"
	"time"
)

func buildCadenceStore() Store {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	s := NewMemoryStore()
	// Host A: very regular — every 60 s
	for i := 0; i < 5; i++ {
		s.Append(Entry{Host: "a", Time: now.Add(time.Duration(i) * 60 * time.Second), Ports: []int{80}})
	}
	// Host B: irregular intervals
	s.Append(Entry{Host: "b", Time: now, Ports: []int{443}})
	s.Append(Entry{Host: "b", Time: now.Add(10 * time.Second), Ports: []int{443}})
	s.Append(Entry{Host: "b", Time: now.Add(300 * time.Second), Ports: []int{443}})
	s.Append(Entry{Host: "b", Time: now.Add(310 * time.Second), Ports: []int{443}})
	return s
}

func TestCadence_Basic(t *testing.T) {
	s := buildCadenceStore()
	results := Cadence(s, 2, 5.0)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestCadence_RegularHost(t *testing.T) {
	s := buildCadenceStore()
	results := Cadence(s, 2, 5.0)
	var ra *CadenceResult
	for i := range results {
		if results[i].Host == "a" {
			ra = &results[i]
		}
	}
	if ra == nil {
		t.Fatal("host a not found")
	}
	if !ra.Regular {
		t.Errorf("expected host a to be regular, jitter=%.2f", ra.Jitter)
	}
	if ra.AvgInterval != 60*time.Second {
		t.Errorf("expected avg 60s, got %v", ra.AvgInterval)
	}
}

func TestCadence_IrregularHost(t *testing.T) {
	s := buildCadenceStore()
	results := Cadence(s, 2, 5.0)
	var rb *CadenceResult
	for i := range results {
		if results[i].Host == "b" {
			rb = &results[i]
		}
	}
	if rb == nil {
		t.Fatal("host b not found")
	}
	if rb.Regular {
		t.Errorf("expected host b to be irregular, jitter=%.2f", rb.Jitter)
	}
}

func TestCadence_MinScansFilters(t *testing.T) {
	s := buildCadenceStore()
	results := Cadence(s, 10, 5.0)
	if len(results) != 0 {
		t.Errorf("expected 0 results with high minScans, got %d", len(results))
	}
}

func TestCadence_EmptyStore(t *testing.T) {
	s := NewMemoryStore()
	results := Cadence(s, 2, 5.0)
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}
