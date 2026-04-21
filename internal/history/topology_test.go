package history

import (
	"testing"
	"time"
)

func buildTopologyStore() Store {
	now := time.Now()
	s := NewMemoryStore()
	// host-a: ports 80, 443
	_ = s.Append(Entry{Host: "host-a", Ports: []int{80, 443}, ScannedAt: now.Add(-2 * time.Minute)})
	_ = s.Append(Entry{Host: "host-a", Ports: []int{80, 443}, ScannedAt: now.Add(-1 * time.Minute)})
	// host-b: ports 443, 8080
	_ = s.Append(Entry{Host: "host-b", Ports: []int{443, 8080}, ScannedAt: now.Add(-2 * time.Minute)})
	_ = s.Append(Entry{Host: "host-b", Ports: []int{443, 8080}, ScannedAt: now.Add(-1 * time.Minute)})
	// host-c: ports 22 only (no shared with others)
	_ = s.Append(Entry{Host: "host-c", Ports: []int{22}, ScannedAt: now.Add(-1 * time.Minute)})
	return s
}

func TestTopology_Nodes(t *testing.T) {
	res := Topology(buildTopologyStore())
	if len(res.Nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(res.Nodes))
	}
	if res.Nodes[0].Host != "host-a" {
		t.Errorf("expected host-a first, got %s", res.Nodes[0].Host)
	}
}

func TestTopology_EdgeSharedPort(t *testing.T) {
	res := Topology(buildTopologyStore())
	// host-a and host-b share port 443
	if len(res.Edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(res.Edges))
	}
	e := res.Edges[0]
	if e.HostA != "host-a" || e.HostB != "host-b" {
		t.Errorf("unexpected edge hosts: %s <-> %s", e.HostA, e.HostB)
	}
	if len(e.SharedPorts) != 1 || e.SharedPorts[0] != 443 {
		t.Errorf("expected shared port 443, got %v", e.SharedPorts)
	}
}

func TestTopology_EmptyStore(t *testing.T) {
	res := Topology(NewMemoryStore())
	if len(res.Nodes) != 0 || len(res.Edges) != 0 {
		t.Errorf("expected empty topology for empty store")
	}
}

func TestTopology_NoSharedPorts(t *testing.T) {
	now := time.Now()
	s := NewMemoryStore()
	_ = s.Append(Entry{Host: "a", Ports: []int{80}, ScannedAt: now})
	_ = s.Append(Entry{Host: "b", Ports: []int{443}, ScannedAt: now})
	res := Topology(s)
	if len(res.Edges) != 0 {
		t.Errorf("expected no edges when no shared ports, got %d", len(res.Edges))
	}
}
