package history

import (
	"testing"
	"time"
)

func buildHeatmapStore() Store {
	base := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	return Store{
		"host-a": {
			{Time: base, Ports: []int{80, 443}},
			{Time: base.Add(2 * time.Hour), Ports: []int{80, 8080}},
			{Time: base.Add(26 * time.Hour), Ports: []int{22}},
		},
		"host-b": {
			{Time: base, Ports: []int{3306}},
		},
	}
}

func TestHeatmap_Basic(t *testing.T) {
	store := buildHeatmapStore()
	entries := Heatmap(store)
	if len(entries) == 0 {
		t.Fatal("expected heatmap entries")
	}
}

func TestHeatmap_MergesPortsPerDay(t *testing.T) {
	store := buildHeatmapStore()
	entries := Heatmap(store)

	for _, e := range entries {
		if e.Host == "host-a" && e.Bucket == "2024-01-15" {
			// ports 80, 443, 8080 => 3 distinct
			if e.Count != 3 {
				t.Errorf("expected 3 distinct ports, got %d", e.Count)
			}
			return
		}
	}
	t.Error("expected entry for host-a on 2024-01-15")
}

func TestHeatmap_SeparateDays(t *testing.T) {
	store := buildHeatmapStore()
	entries := Heatmap(store)

	buckets := map[string]int{}
	for _, e := range entries {
		if e.Host == "host-a" {
			buckets[e.Bucket] = e.Count
		}
	}
	if _, ok := buckets["2024-01-16"]; !ok {
		t.Error("expected entry for host-a on 2024-01-16")
	}
	if buckets["2024-01-16"] != 1 {
		t.Errorf("expected 1 port on 2024-01-16, got %d", buckets["2024-01-16"])
	}
}

func TestHeatmap_EmptyStore(t *testing.T) {
	entries := Heatmap(Store{})
	if len(entries) != 0 {
		t.Errorf("expected no entries, got %d", len(entries))
	}
}
