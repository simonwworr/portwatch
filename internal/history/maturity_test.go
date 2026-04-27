package history

import (
	"testing"
	"time"
)

func buildMaturityStore() Store {
	now := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	s := NewMemoryStore()

	// hostA: 10 scans over 30 days, ports 80 and 443 always open → high maturity
	for i := 0; i < 10; i++ {
		s.Append(Entry{
			Host:  "hostA",
			Time:  now.AddDate(0, 0, -30+i*3),
			Ports: []int{80, 443},
		})
	}

	// hostB: 5 scans, port set changes every scan → low maturity
	for i := 0; i < 5; i++ {
		s.Append(Entry{
			Host:  "hostB",
			Time:  now.AddDate(0, 0, -10+i*2),
			Ports: []int{8000 + i},
		})
	}

	// hostC: only 1 scan → filtered by minScans
	s.Append(Entry{
		Host:  "hostC",
		Time:  now.AddDate(0, 0, -5),
		Ports: []int{22},
	})

	return s
}

func TestMaturity_Basic(t *testing.T) {
	now := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	results := Maturity(buildMaturityStore(), 2, now)
	if len(results) != 2 {
		t.Fatalf("expected 2 results (hostC filtered), got %d", len(results))
	}
}

func TestMaturity_HostAScoresHighest(t *testing.T) {
	now := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	results := Maturity(buildMaturityStore(), 2, now)
	if results[0].Host != "hostA" {
		t.Errorf("expected hostA to score highest, got %s", results[0].Host)
	}
}

func TestMaturity_StablePorts(t *testing.T) {
	now := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	results := Maturity(buildMaturityStore(), 2, now)
	var hostA MaturityResult
	for _, r := range results {
		if r.Host == "hostA" {
			hostA = r
		}
	}
	if len(hostA.StablePorts) != 2 {
		t.Errorf("expected 2 stable ports for hostA, got %v", hostA.StablePorts)
	}
}

func TestMaturity_MinScansFilters(t *testing.T) {
	now := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	results := Maturity(buildMaturityStore(), 6, now)
	if len(results) != 1 || results[0].Host != "hostA" {
		t.Errorf("expected only hostA with minScans=6, got %+v", results)
	}
}

func TestMaturity_EmptyStore(t *testing.T) {
	now := time.Now()
	results := Maturity(NewMemoryStore(), 1, now)
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}
