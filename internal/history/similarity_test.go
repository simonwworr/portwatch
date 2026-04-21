package history

import (
	"testing"
	"time"
)

func buildSimilarityStore() Store {
	now := time.Now()
	s := NewMemoryStore()
	s.Append("host-a", Entry{Timestamp: now, Ports: []int{80, 443, 8080}})
	s.Append("host-b", Entry{Timestamp: now, Ports: []int{80, 443, 9090}})
	s.Append("host-c", Entry{Timestamp: now, Ports: []int{22, 3306}})
	return s
}

func TestSimilarity_Basic(t *testing.T) {
	s := buildSimilarityStore()
	results := Similarity(s, 0.0)
	if len(results) == 0 {
		t.Fatal("expected at least one result")
	}
}

func TestSimilarity_HostAandB(t *testing.T) {
	s := buildSimilarityStore()
	results := Similarity(s, 0.0)
	var found *SimilarityResult
	for i := range results {
		r := &results[i]
		if (r.HostA == "host-a" && r.HostB == "host-b") || (r.HostA == "host-b" && r.HostB == "host-a") {
			found = r
			break
		}
	}
	if found == nil {
		t.Fatal("expected pair host-a / host-b")
	}
	// shared: 80,443 union: 80,443,8080,9090 => 2/4 = 0.5
	if found.Similarity != 0.5 {
		t.Errorf("expected 0.5, got %v", found.Similarity)
	}
	if len(found.Shared) != 2 {
		t.Errorf("expected 2 shared ports, got %d", len(found.Shared))
	}
}

func TestSimilarity_MinScore(t *testing.T) {
	s := buildSimilarityStore()
	results := Similarity(s, 0.9)
	for _, r := range results {
		if r.Similarity < 0.9 {
			t.Errorf("result below minScore: %v", r.Similarity)
		}
	}
}

func TestSimilarity_EmptyStore(t *testing.T) {
	s := NewMemoryStore()
	results := Similarity(s, 0.0)
	if len(results) != 0 {
		t.Errorf("expected no results for empty store")
	}
}

func TestSimilarity_SortedDescending(t *testing.T) {
	now := time.Now()
	s := NewMemoryStore()
	s.Append("h1", Entry{Timestamp: now, Ports: []int{80, 443}})
	s.Append("h2", Entry{Timestamp: now, Ports: []int{80, 443}})
	s.Append("h3", Entry{Timestamp: now, Ports: []int{80, 9999}})
	results := Similarity(s, 0.0)
	for i := 1; i < len(results); i++ {
		if results[i].Similarity > results[i-1].Similarity {
			t.Errorf("results not sorted descending at index %d", i)
		}
	}
}
