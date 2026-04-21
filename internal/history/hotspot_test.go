package history

import (
	"testing"
	"time"
)

func buildHotspotStore() Store {
	now := time.Now()
	ms := NewMemoryStore()
	// port 80 open on both hostA and hostB in every scan
	for i := 0; i < 3; i++ {
		_ = ms.Append(Entry{Host: "hostA", Ports: []int{80, 443}, Time: now.Add(time.Duration(i) * time.Hour)})
		_ = ms.Append(Entry{Host: "hostB", Ports: []int{80, 8080}, Time: now.Add(time.Duration(i) * time.Hour)})
	}
	// port 9090 only on hostC (single host)
	_ = ms.Append(Entry{Host: "hostC", Ports: []int{9090}, Time: now})
	return ms
}

func TestHotspot_FindsMultiHostPort(t *testing.T) {
	store := buildHotspotStore()
	results := Hotspot(store, 2)
	if len(results) == 0 {
		t.Fatal("expected at least one hotspot")
	}
	// port 80 should be top result (on both hostA and hostB)
	if results[0].Port != 80 {
		t.Errorf("expected port 80 as top hotspot, got %d", results[0].Port)
	}
	if results[0].HostCount != 2 {
		t.Errorf("expected host_count=2, got %d", results[0].HostCount)
	}
}

func TestHotspot_FrequencyIsOne_WhenAlwaysOpen(t *testing.T) {
	store := buildHotspotStore()
	results := Hotspot(store, 2)
	for _, r := range results {
		if r.Port == 80 {
			if r.Frequency < 0.99 {
				t.Errorf("expected frequency ~1.0 for port 80, got %.2f", r.Frequency)
			}
			return
		}
	}
	t.Error("port 80 not found in results")
}

func TestHotspot_MinHostsFilters(t *testing.T) {
	store := buildHotspotStore()
	// minHosts=3 means no port qualifies (max is 2)
	results := Hotspot(store, 3)
	if len(results) != 0 {
		t.Errorf("expected no results with minHosts=3, got %d", len(results))
	}
}

func TestHotspot_EmptyStore(t *testing.T) {
	ms := NewMemoryStore()
	results := Hotspot(ms, 1)
	if results != nil {
		t.Errorf("expected nil for empty store, got %v", results)
	}
}

func TestHotspot_HostsAreSorted(t *testing.T) {
	store := buildHotspotStore()
	results := Hotspot(store, 2)
	for _, r := range results {
		for i := 1; i < len(r.Hosts); i++ {
			if r.Hosts[i] < r.Hosts[i-1] {
				t.Errorf("hosts not sorted for port %d", r.Port)
			}
		}
	}
}
