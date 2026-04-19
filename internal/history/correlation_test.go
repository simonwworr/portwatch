package history

import (
	"testing"
	"time"
)

func buildCorrelationStore() Store {
	now := time.Now()
	return Store{
		Entries: []Entry{
			{Host: "host1", Ports: []int{80, 443, 8080}, Time: now},
			{Host: "host1", Ports: []int{80, 443}, Time: now.Add(time.Hour)},
			{Host: "host1", Ports: []int{80, 443, 22}, Time: now.Add(2 * time.Hour)},
			{Host: "host2", Ports: []int{22, 3306}, Time: now},
			{Host: "host2", Ports: []int{22, 3306}, Time: now.Add(time.Hour)},
		},
	}
}

func TestCorrelate_Basic(t *testing.T) {
	store := buildCorrelationStore()
	results := Correlate(store, 2)
	if len(results) == 0 {
		t.Fatal("expected results")
	}
	var host1 *CorrelationResult
	for i := range results {
		if results[i].Host == "host1" {
			host1 = &results[i]
		}
	}
	if host1 == nil {
		t.Fatal("expected host1 in results")
	}
	if host1.Pairs[0].PortA != 80 || host1.Pairs[0].PortB != 443 {
		t.Errorf("expected 80/443 as top pair, got %d/%d", host1.Pairs[0].PortA, host1.Pairs[0].PortB)
	}
	if host1.Pairs[0].Count != 3 {
		t.Errorf("expected count 3, got %d", host1.Pairs[0].Count)
	}
}

func TestCorrelate_MinCount(t *testing.T) {
	store := buildCorrelationStore()
	results := Correlate(store, 3)
	for _, r := range results {
		for _, p := range r.Pairs {
			if p.Count < 3 {
				t.Errorf("pair %d/%d has count %d below minCount", p.PortA, p.PortB, p.Count)
			}
		}
	}
}

func TestCorrelate_EmptyStore(t *testing.T) {
	results := Correlate(Store{}, 1)
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}

func TestCorrelate_Host2(t *testing.T) {
	store := buildCorrelationStore()
	results := Correlate(store, 2)
	var host2 *CorrelationResult
	for i := range results {
		if results[i].Host == "host2" {
			host2 = &results[i]
		}
	}
	if host2 == nil {
		t.Fatal("expected host2")
	}
	if len(host2.Pairs) != 1 || host2.Pairs[0].PortA != 22 || host2.Pairs[0].PortB != 3306 {
		t.Errorf("unexpected pairs for host2: %+v", host2.Pairs)
	}
}
