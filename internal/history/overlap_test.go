package history

import (
	"testing"
	"time"
)

func buildOverlapStore(t *testing.T) Store {
	t.Helper()
	now := time.Now()
	store := NewMemoryStore()

	// hostA: ports 80, 443, 8080
	store.Append("hostA", Entry{Timestamp: now, Ports: []int{80, 443, 8080}})
	// hostB: ports 80, 443, 9090
	store.Append("hostB", Entry{Timestamp: now, Ports: []int{80, 443, 9090}})
	// hostC: ports 22, 3306
	store.Append("hostC", Entry{Timestamp: now, Ports: []int{22, 3306}})
	return store
}

func TestOverlap_SharedPorts(t *testing.T) {
	store := buildOverlapStore(t)
	results := Overlap(store)

	if len(results) == 0 {
		t.Fatal("expected overlap results, got none")
	}

	// hostA & hostB share 80 and 443 — highest Jaccard
	top := results[0]
	if top.HostA != "hostA" || top.HostB != "hostB" {
		t.Errorf("expected hostA/hostB as top pair, got %s/%s", top.HostA, top.HostB)
	}
	if len(top.SharedPorts) != 2 {
		t.Errorf("expected 2 shared ports, got %d", len(top.SharedPorts))
	}
	if top.SharedPorts[0] != 80 || top.SharedPorts[1] != 443 {
		t.Errorf("unexpected shared ports: %v", top.SharedPorts)
	}
}

func TestOverlap_JaccardScore(t *testing.T) {
	store := buildOverlapStore(t)
	results := Overlap(store)

	top := results[0]
	// |{80,443}| / |{80,443,8080,9090}| = 2/4 = 0.5
	if top.JaccardScore < 0.49 || top.JaccardScore > 0.51 {
		t.Errorf("expected Jaccard ~0.5, got %.4f", top.JaccardScore)
	}
}

func TestOverlap_NoSharedPorts(t *testing.T) {
	store := buildOverlapStore(t)
	results := Overlap(store)

	// Find hostA/hostC or hostB/hostC pair — no shared ports
	var found *OverlapResult
	for i := range results {
		r := &results[i]
		if (r.HostA == "hostA" && r.HostB == "hostC") ||
			(r.HostA == "hostB" && r.HostB == "hostC") {
			found = r
			break
		}
	}
	if found == nil {
		t.Fatal("expected a result involving hostC")
	}
	if len(found.SharedPorts) != 0 {
		t.Errorf("expected no shared ports with hostC, got %v", found.SharedPorts)
	}
	if found.JaccardScore != 0 {
		t.Errorf("expected Jaccard 0, got %.4f", found.JaccardScore)
	}
}

func TestOverlap_EmptyStore(t *testing.T) {
	store := NewMemoryStore()
	results := Overlap(store)
	if results != nil {
		t.Errorf("expected nil for empty store, got %v", results)
	}
}

func TestOverlap_SingleHost(t *testing.T) {
	store := NewMemoryStore()
	store.Append("hostA", Entry{Timestamp: time.Now(), Ports: []int{80}})
	results := Overlap(store)
	if len(results) != 0 {
		t.Errorf("expected no results for single host, got %d", len(results))
	}
}
