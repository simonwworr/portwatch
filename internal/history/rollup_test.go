package history

import (
	"testing"
	"time"
)

func buildRollupStore() Store {
	now := time.Now()
	return Store{
		"host-a": {
			{Host: "host-a", Ports: []int{80, 443}, ScannedAt: now.Add(-2 * time.Hour)},
			{Host: "host-a", Ports: []int{80, 8080}, ScannedAt: now.Add(-1 * time.Hour)},
			{Host: "host-a", Ports: []int{80}, ScannedAt: now},
		},
		"host-b": {
			{Host: "host-b", Ports: []int{22, 3306}, ScannedAt: now},
		},
	}
}

func TestRollup_Basic(t *testing.T) {
	store := buildRollupStore()
	result := Rollup(store)
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result[0].Host != "host-a" {
		t.Errorf("expected host-a first, got %s", result[0].Host)
	}
}

func TestRollup_ScanCount(t *testing.T) {
	store := buildRollupStore()
	result := Rollup(store)
	for _, r := range result {
		if r.Host == "host-a" && r.ScanCount != 3 {
			t.Errorf("expected 3 scans for host-a, got %d", r.ScanCount)
		}
		if r.Host == "host-b" && r.ScanCount != 1 {
			t.Errorf("expected 1 scan for host-b, got %d", r.ScanCount)
		}
	}
}

func TestRollup_UniquePorts(t *testing.T) {
	store := buildRollupStore()
	result := Rollup(store)
	for _, r := range result {
		if r.Host == "host-a" {
			// ports seen: 80, 443, 8080
			if len(r.UniquePorts) != 3 {
				t.Errorf("expected 3 unique ports, got %d", len(r.UniquePorts))
			}
		}
	}
}

func TestRollup_MinMax(t *testing.T) {
	store := buildRollupStore()
	result := Rollup(store)
	for _, r := range result {
		if r.Host == "host-a" {
			if r.MinOpen != 1 {
				t.Errorf("expected MinOpen=1, got %d", r.MinOpen)
			}
			if r.MaxOpen != 2 {
				t.Errorf("expected MaxOpen=2, got %d", r.MaxOpen)
			}
		}
	}
}

func TestRollup_EmptyStore(t *testing.T) {
	result := Rollup(Store{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d", len(result))
	}
}
