package history

import (
	"testing"
	"time"
)

func buildMomentumStore(t *testing.T) Store {
	t.Helper()
	store := NewMemoryStore()
	now := time.Now()

	// host-a: consistently opening ports over 4 scans
	store.Append("host-a", Entry{Time: now.Add(-3 * 24 * time.Hour), Ports: []int{80}})
	store.Append("host-a", Entry{Time: now.Add(-2 * 24 * time.Hour), Ports: []int{80, 443}})
	store.Append("host-a", Entry{Time: now.Add(-1 * 24 * time.Hour), Ports: []int{80, 443, 8080}})
	store.Append("host-a", Entry{Time: now, Ports: []int{80, 443, 8080, 9090}})

	// host-b: consistently closing ports
	store.Append("host-b", Entry{Time: now.Add(-3 * 24 * time.Hour), Ports: []int{80, 443, 8080, 9090}})
	store.Append("host-b", Entry{Time: now.Add(-2 * 24 * time.Hour), Ports: []int{80, 443, 8080}})
	store.Append("host-b", Entry{Time: now.Add(-1 * 24 * time.Hour), Ports: []int{80, 443}})
	store.Append("host-b", Entry{Time: now, Ports: []int{80}})

	// host-c: stable, no change
	store.Append("host-c", Entry{Time: now.Add(-2 * 24 * time.Hour), Ports: []int{80, 443}})
	store.Append("host-c", Entry{Time: now.Add(-1 * 24 * time.Hour), Ports: []int{80, 443}})
	store.Append("host-c", Entry{Time: now, Ports: []int{80, 443}})

	return store
}

func TestMomentum_Basic(t *testing.T) {
	store := buildMomentumStore(t)
	results := Momentum(store, 2, 0.5)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
}

func TestMomentum_OpeningHostScoresPositive(t *testing.T) {
	store := buildMomentumStore(t)
	results := Momentum(store, 2, 0.5)
	var score float64
	for _, r := range results {
		if r.Host == "host-a" {
			score = r.Score
		}
	}
	if score <= 0 {
		t.Errorf("host-a should have positive momentum, got %f", score)
	}
}

func TestMomentum_ClosingHostScoresNegative(t *testing.T) {
	store := buildMomentumStore(t)
	results := Momentum(store, 2, 0.5)
	var score float64
	for _, r := range results {
		if r.Host == "host-b" {
			score = r.Score
		}
	}
	if score >= 0 {
		t.Errorf("host-b should have negative momentum, got %f", score)
	}
}

func TestMomentum_StableHostScoresZero(t *testing.T) {
	store := buildMomentumStore(t)
	results := Momentum(store, 2, 0.5)
	var score float64
	for _, r := range results {
		if r.Host == "host-c" {
			score = r.Score
		}
	}
	if score != 0.0 {
		t.Errorf("host-c should have zero momentum, got %f", score)
	}
}

func TestMomentum_MinScansFilters(t *testing.T) {
	store := NewMemoryStore()
	now := time.Now()
	store.Append("host-x", Entry{Time: now, Ports: []int{80}})
	results := Momentum(store, 3, 0.5)
	if len(results) != 0 {
		t.Errorf("expected 0 results with minScans=3, got %d", len(results))
	}
}

func TestMomentum_EmptyStore(t *testing.T) {
	store := NewMemoryStore()
	results := Momentum(store, 2, 0.5)
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}
