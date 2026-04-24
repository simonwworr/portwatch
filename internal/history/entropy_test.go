package history

import (
	"math"
	"testing"
	"time"
)

func buildEntropyStore() Store {
	now := time.Now()
	m := NewMemoryStore()
	// host-a: ports vary a lot — high entropy
	for i := 0; i < 6; i++ {
		m.Append("host-a", Entry{Timestamp: now.Add(time.Duration(i) * time.Hour), Ports: []int{80, 443, 8080 + i}})
	}
	// host-b: always same ports — low entropy
	for i := 0; i < 6; i++ {
		m.Append("host-b", Entry{Timestamp: now.Add(time.Duration(i) * time.Hour), Ports: []int{80, 443}})
	}
	// host-c: too few scans
	m.Append("host-c", Entry{Timestamp: now, Ports: []int{22}})
	return m
}

func TestEntropy_Basic(t *testing.T) {
	store := buildEntropyStore()
	results := Entropy(store, 2)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestEntropy_HigherForVariedHost(t *testing.T) {
	store := buildEntropyStore()
	results := Entropy(store, 2)

	var ea, eb float64
	for _, r := range results {
		if r.Host == "host-a" {
			ea = r.Entropy
		}
		if r.Host == "host-b" {
			eb = r.Entropy
		}
	}
	if ea <= eb {
		t.Errorf("expected host-a entropy (%.3f) > host-b entropy (%.3f)", ea, eb)
	}
}

func TestEntropy_MinScansFilters(t *testing.T) {
	store := buildEntropyStore()
	results := Entropy(store, 10)
	if len(results) != 0 {
		t.Errorf("expected 0 results with high minScans, got %d", len(results))
	}
}

func TestEntropy_EmptyStore(t *testing.T) {
	m := NewMemoryStore()
	results := Entropy(m, 1)
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}

func TestEntropy_UniformDistribution(t *testing.T) {
	// Two ports, each appearing exactly half the time => entropy = 1.0
	now := time.Now()
	m := NewMemoryStore()
	for i := 0; i < 4; i++ {
		ports := []int{80}
		if i%2 == 0 {
			ports = []int{443}
		}
		m.Append("host-x", Entry{Timestamp: now.Add(time.Duration(i) * time.Hour), Ports: ports})
	}
	results := Entropy(m, 1)
	if len(results) != 1 {
		t.Fatalf("expected 1 result")
	}
	if math.Abs(results[0].Entropy-1.0) > 0.01 {
		t.Errorf("expected entropy ~1.0, got %.3f", results[0].Entropy)
	}
}
