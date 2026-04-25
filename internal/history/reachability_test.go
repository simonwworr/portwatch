package history

import (
	"testing"
	"time"
)

func buildReachabilityStore() Store {
	now := time.Now()
	s := NewMemoryStore()
	// hostA: 3 scans, all active
	s.Append("hostA", Entry{Timestamp: now.Add(-2 * time.Hour), Ports: []int{80, 443}})
	s.Append("hostA", Entry{Timestamp: now.Add(-1 * time.Hour), Ports: []int{80}})
	s.Append("hostA", Entry{Timestamp: now, Ports: []int{80, 443, 8080}})
	// hostB: 3 scans, one inactive
	s.Append("hostB", Entry{Timestamp: now.Add(-2 * time.Hour), Ports: []int{22}})
	s.Append("hostB", Entry{Timestamp: now.Add(-1 * time.Hour), Ports: []int{}})
	s.Append("hostB", Entry{Timestamp: now, Ports: []int{22}})
	// hostC: 1 scan — below minScans=2
	s.Append("hostC", Entry{Timestamp: now, Ports: []int{80}})
	return s
}

func TestReachability_Basic(t *testing.T) {
	s := buildReachabilityStore()
	res := Reachability(s, 2)
	if len(res) != 2 {
		t.Fatalf("expected 2 results, got %d", len(res))
	}
}

func TestReachability_HostAFullUptime(t *testing.T) {
	s := buildReachabilityStore()
	res := Reachability(s, 2)
	// hostA should be first (uptime=1.0)
	if res[0].Host != "hostA" {
		t.Fatalf("expected hostA first, got %s", res[0].Host)
	}
	if res[0].Uptime != 1.0 {
		t.Errorf("expected uptime 1.0, got %f", res[0].Uptime)
	}
}

func TestReachability_HostBPartialUptime(t *testing.T) {
	s := buildReachabilityStore()
	res := Reachability(s, 2)
	var hostB *ReachabilityResult
	for i := range res {
		if res[i].Host == "hostB" {
			hostB = &res[i]
		}
	}
	if hostB == nil {
		t.Fatal("hostB not found")
	}
	expected := 2.0 / 3.0
	if hostB.Uptime < expected-0.001 || hostB.Uptime > expected+0.001 {
		t.Errorf("expected uptime ~%.4f, got %.4f", expected, hostB.Uptime)
	}
}

func TestReachability_MinScansFilters(t *testing.T) {
	s := buildReachabilityStore()
	res := Reachability(s, 4)
	if len(res) != 0 {
		t.Errorf("expected 0 results with minScans=4, got %d", len(res))
	}
}

func TestReachability_EmptyStore(t *testing.T) {
	s := NewMemoryStore()
	res := Reachability(s, 1)
	if len(res) != 0 {
		t.Errorf("expected empty result, got %d", len(res))
	}
}

func TestReachability_MaxMinPorts(t *testing.T) {
	s := buildReachabilityStore()
	res := Reachability(s, 2)
	var hostA *ReachabilityResult
	for i := range res {
		if res[i].Host == "hostA" {
			hostA = &res[i]
		}
	}
	if hostA == nil {
		t.Fatal("hostA not found")
	}
	if hostA.MaxPorts != 3 {
		t.Errorf("expected MaxPorts=3, got %d", hostA.MaxPorts)
	}
	if hostA.MinPorts != 1 {
		t.Errorf("expected MinPorts=1, got %d", hostA.MinPorts)
	}
}
