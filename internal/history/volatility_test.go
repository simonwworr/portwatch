package history

import (
	"testing"
	"time"
)

func buildVolatilityStore() Store {
	now := time.Now()
	s := NewMemoryStore()

	// host-a: changes every scan → high volatility
	s.Append("host-a", Entry{Timestamp: now.Add(-4 * 24 * time.Hour), Ports: []int{80, 443}})
	s.Append("host-a", Entry{Timestamp: now.Add(-3 * 24 * time.Hour), Ports: []int{80, 8080}})
	s.Append("host-a", Entry{Timestamp: now.Add(-2 * 24 * time.Hour), Ports: []int{443, 9090}})
	s.Append("host-a", Entry{Timestamp: now.Add(-1 * 24 * time.Hour), Ports: []int{80, 443}})

	// host-b: never changes → zero volatility
	s.Append("host-b", Entry{Timestamp: now.Add(-3 * 24 * time.Hour), Ports: []int{80, 443}})
	s.Append("host-b", Entry{Timestamp: now.Add(-2 * 24 * time.Hour), Ports: []int{80, 443}})
	s.Append("host-b", Entry{Timestamp: now.Add(-1 * 24 * time.Hour), Ports: []int{80, 443}})

	return s
}

func TestVolatility_Basic(t *testing.T) {
	store := buildVolatilityStore()
	results := Volatility(store, 2)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestVolatility_HighForChangingHost(t *testing.T) {
	store := buildVolatilityStore()
	results := Volatility(store, 2)

	var hostA VolatilityResult
	for _, r := range results {
		if r.Host == "host-a" {
			hostA = r
		}
	}

	if hostA.Volatility <= 0 {
		t.Errorf("expected positive volatility for host-a, got %f", hostA.Volatility)
	}
	if hostA.Changes == 0 {
		t.Errorf("expected changes > 0 for host-a")
	}
}

func TestVolatility_ZeroForStableHost(t *testing.T) {
	store := buildVolatilityStore()
	results := Volatility(store, 2)

	for _, r := range results {
		if r.Host == "host-b" {
			if r.Volatility != 0 {
				t.Errorf("expected zero volatility for host-b, got %f", r.Volatility)
			}
			if r.Changes != 0 {
				t.Errorf("expected zero changes for host-b, got %d", r.Changes)
			}
		}
	}
}

func TestVolatility_MinScansFilters(t *testing.T) {
	now := time.Now()
	s := NewMemoryStore()
	s.Append("host-c", Entry{Timestamp: now, Ports: []int{80}})

	results := Volatility(s, 3)
	if len(results) != 0 {
		t.Errorf("expected no results when below minScans, got %d", len(results))
	}
}

func TestVolatility_EmptyStore(t *testing.T) {
	s := NewMemoryStore()
	results := Volatility(s, 2)
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}
