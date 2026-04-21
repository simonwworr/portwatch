package history

import "sort"

// TopologyNode represents a host and its current open ports.
type TopologyNode struct {
	Host  string
	Ports []int
}

// TopologyEdge represents two hosts that share at least one open port.
type TopologyEdge struct {
	HostA        string
	HostB        string
	SharedPorts  []int
}

// TopologyResult holds the full topology map.
type TopologyResult struct {
	Nodes []TopologyNode
	Edges []TopologyEdge
}

// Topology builds a port-sharing graph from the latest snapshot of each host.
// Two hosts are connected by an edge when they share one or more open ports.
func Topology(store Store) TopologyResult {
	snaps := AllSnapshots(store)

	nodes := make([]TopologyNode, 0, len(snaps))
	for _, s := range snaps {
		ports := make([]int, len(s.Ports))
		copy(ports, s.Ports)
		sort.Ints(ports)
		nodes = append(nodes, TopologyNode{Host: s.Host, Ports: ports})
	}
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].Host < nodes[j].Host })

	var edges []TopologyEdge
	for i := 0; i < len(nodes); i++ {
		for j := i + 1; j < len(nodes); j++ {
			shared := sharedPorts(nodes[i].Ports, nodes[j].Ports)
			if len(shared) > 0 {
				edges = append(edges, TopologyEdge{
					HostA:       nodes[i].Host,
					HostB:       nodes[j].Host,
					SharedPorts: shared,
				})
			}
		}
	}

	return TopologyResult{Nodes: nodes, Edges: edges}
}

func sharedPorts(a, b []int) []int {
	set := make(map[int]struct{}, len(a))
	for _, p := range a {
		set[p] = struct{}{}
	}
	var shared []int
	for _, p := range b {
		if _, ok := set[p]; ok {
			shared = append(shared, p)
		}
	}
	sort.Ints(shared)
	return shared
}
