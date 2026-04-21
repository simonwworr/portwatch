package history

import (
	"math"
	"sort"
)

// PageRankResult holds the computed importance score for a host.
type PageRankResult struct {
	Host  string
	Score float64
}

// PageRank computes a simplified PageRank-style importance score for each host
// based on shared open ports (edges in the topology graph). Hosts with more
// shared-port connections to other highly-connected hosts score higher.
//
// iterations controls convergence (typically 20–50 is sufficient).
// dampingFactor is usually 0.85.
func PageRank(store Store, iterations int, dampingFactor float64) []PageRankResult {
	topo := Topology(store)
	if len(topo.Nodes) == 0 {
		return nil
	}

	// Build adjacency: host -> set of connected hosts
	adj := make(map[string][]string)
	for _, node := range topo.Nodes {
		adj[node.Host] = []string{}
	}
	for _, edge := range topo.Edges {
		adj[edge.HostA] = append(adj[edge.HostA], edge.HostB)
		adj[edge.HostB] = append(adj[edge.HostB], edge.HostA)
	}

	n := float64(len(topo.Nodes))
	scores := make(map[string]float64)
	for _, node := range topo.Nodes {
		scores[node.Host] = 1.0 / n
	}

	for i := 0; i < iterations; i++ {
		next := make(map[string]float64)
		for _, node := range topo.Nodes {
			next[node.Host] = (1 - dampingFactor) / n
		}
		for host, neighbors := range adj {
			if len(neighbors) == 0 {
				continue
			}
			contrib := dampingFactor * scores[host] / float64(len(neighbors))
			for _, nb := range neighbors {
				next[nb] += contrib
			}
		}
		scores = next
	}

	results := make([]PageRankResult, 0, len(scores))
	for host, s := range scores {
		results = append(results, PageRankResult{
			Host:  host,
			Score: math.Round(s*1e6) / 1e6,
		})
	}
	sort.Slice(results, func(i, j int) bool {
		if results[i].Score != results[j].Score {
			return results[i].Score > results[j].Score
		}
		return results[i].Host < results[j].Host
	})
	return results
}
