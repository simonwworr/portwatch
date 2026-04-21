package history

import (
	"testing"
	"time"
)

func buildPageRankStore() Store {
	now := time.Now()
	s := NewMemoryStore()
	// hostA and hostB share port 80 and 443
	// hostA and hostC share port 22
	// hostB and hostC share no ports
	// => hostA should score highest (most connections)
	_ = s.Append("hostA", Entry{Timestamp: now, Ports: []int{22, 80, 443}})
	_ = s.Append("hostB", Entry{Timestamp: now, Ports: []int{80, 443, 8080}})
	_ = s.Append("hostC", Entry{Timestamp: now, Ports: []int{22, 9000}})
	return s
}

func TestPageRank_HostAScoresHighest(t *testing.T) {
	store := buildPageRankStore()
	results := PageRank(store, 30, 0.85)
	if len(results) == 0 {
		t.Fatal("expected results, got none")
	}
	if results[0].Host != "hostA" {
		t.Errorf("expected hostA to rank first, got %s", results[0].Host)
	}
}

func TestPageRank_AllHostsPresent(t *testing.T) {
	store := buildPageRankStore()
	results := PageRank(store, 30, 0.85)
	hosts := map[string]bool{}
	for _, r := range results {
		hosts[r.Host] = true
	}
	for _, h := range []string{"hostA", "hostB", "hostC"} {
		if !hosts[h] {
			t.Errorf("missing host %s in results", h)
		}
	}
}

func TestPageRank_ScoresPositive(t *testing.T) {
	store := buildPageRankStore()
	results := PageRank(store, 30, 0.85)
	for _, r := range results {
		if r.Score <= 0 {
			t.Errorf("host %s has non-positive score %f", r.Host, r.Score)
		}
	}
}

func TestPageRank_EmptyStore(t *testing.T) {
	s := NewMemoryStore()
	results := PageRank(s, 30, 0.85)
	if len(results) != 0 {
		t.Errorf("expected empty results for empty store, got %d", len(results))
	}
}

func TestPageRank_SingleHost(t *testing.T) {
	s := NewMemoryStore()
	_ = s.Append("solo", Entry{Timestamp: time.Now(), Ports: []int{80}})
	results := PageRank(s, 30, 0.85)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Host != "solo" {
		t.Errorf("unexpected host %s", results[0].Host)
	}
}
