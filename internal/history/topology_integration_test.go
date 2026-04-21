package history

import (
	"testing"
	"time"
)

// TestTopology_MultipleSharedPorts verifies that when two hosts share several
// ports, all of them appear in the edge's SharedPorts slice.
func TestTopology_MultipleSharedPorts(t *testing.T) {
	now := time.Now()
	s := NewMemoryStore()
	_ = s.Append(Entry{Host: "x", Ports: []int{22, 80, 443}, ScannedAt: now})
	_ = s.Append(Entry{Host: "y", Ports: []int{80, 443, 8443}, ScannedAt: now})

	res := Topology(s)
	if len(res.Edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(res.Edges))
	}
	if len(res.Edges[0].SharedPorts) != 2 {
		t.Errorf("expected 2 shared ports, got %v", res.Edges[0].SharedPorts)
	}
}

// TestTopology_ThreeHostsFullyConnected checks that three hosts sharing a port
// produce three edges.
func TestTopology_ThreeHostsFullyConnected(t *testing.T) {
	now := time.Now()
	s := NewMemoryStore()
	_ = s.Append(Entry{Host: "a", Ports: []int{443}, ScannedAt: now})
	_ = s.Append(Entry{Host: "b", Ports: []int{443}, ScannedAt: now})
	_ = s.Append(Entry{Host: "c", Ports: []int{443}, ScannedAt: now})

	res := Topology(s)
	if len(res.Edges) != 3 {
		t.Errorf("expected 3 edges for fully-connected trio, got %d", len(res.Edges))
	}
	for _, e := range res.Edges {
		if len(e.SharedPorts) != 1 || e.SharedPorts[0] != 443 {
			t.Errorf("unexpected shared ports in edge %+v", e)
		}
	}
}

// TestTopology_NodePortsAreSorted ensures node port lists come back sorted.
func TestTopology_NodePortsAreSorted(t *testing.T) {
	now := time.Now()
	s := NewMemoryStore()
	_ = s.Append(Entry{Host: "h", Ports: []int{8080, 22, 443, 80}, ScannedAt: now})

	res := Topology(s)
	if len(res.Nodes) != 1 {
		t.Fatalf("expected 1 node")
	}
	p := res.Nodes[0].Ports
	for i := 1; i < len(p); i++ {
		if p[i] < p[i-1] {
			t.Errorf("ports not sorted: %v", p)
		}
	}
}
