package history

import (
	"testing"
	"time"
)

func buildSearchStore() *Store {
	s := &Store{}
	now := time.Now()
	s.Append("host-a", Entry{Timestamp: now.Add(-2 * time.Hour), Ports: []int{80, 443}})
	s.Append("host-a", Entry{Timestamp: now.Add(-1 * time.Hour), Ports: []int{80, 8080}})
	s.Append("host-b", Entry{Timestamp: now.Add(-30 * time.Minute), Ports: []int{22, 443}})
	return s
}

func TestSearch_ByHost(t *testing.T) {
	s := buildSearchStore()
	results := Search(s, SearchQuery{Host: "host-a"})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Host != "host-a" {
			t.Errorf("unexpected host %s", r.Host)
		}
	}
}

func TestSearch_ByPort(t *testing.T) {
	s := buildSearchStore()
	results := Search(s, SearchQuery{Port: 443})
	if len(results) != 2 {
		t.Fatalf("expected 2 results for port 443, got %d", len(results))
	}
}

func TestSearch_BySince(t *testing.T) {
	s := buildSearchStore()
	results := Search(s, SearchQuery{Since: time.Now().Add(-90 * time.Minute)})
	if len(results) != 2 {
		t.Fatalf("expected 2 results since 90m ago, got %d", len(results))
	}
}

func TestSearch_NoMatch(t *testing.T) {
	s := buildSearchStore()
	results := Search(s, SearchQuery{Port: 9999})
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestSearch_AllEntries(t *testing.T) {
	s := buildSearchStore()
	results := Search(s, SearchQuery{})
	if len(results) != 3 {
		t.Fatalf("expected 3 total results, got %d", len(results))
	}
}
