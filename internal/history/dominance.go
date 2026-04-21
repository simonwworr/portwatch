package history

import "sort"

// DominanceResult holds the dominance score for a host.
// A host is "dominant" when its open ports are frequently a superset
// of other hosts' open ports.
type DominanceResult struct {
	Host  string
	Score float64 // fraction of other hosts whose latest ports are a subset
	Peers int     // number of hosts compared against
}

// Dominance computes, for each host, what fraction of other hosts have
// a latest port-set that is a subset of that host's latest port-set.
// Only hosts with at least minPorts open ports are included as candidates.
func Dominance(store Store, minPorts int) []DominanceResult {
	hosts := store.Hosts()
	if len(hosts) == 0 {
		return nil
	}

	// Build latest snapshot per host.
	latest := make(map[string]map[int]struct{}, len(hosts))
	for _, h := range hosts {
		entries := store.ForHost(h)
		if len(entries) == 0 {
			continue
		}
		last := entries[len(entries)-1]
		set := make(map[int]struct{}, len(last.Ports))
		for _, p := range last.Ports {
			set[p] = struct{}{}
		}
		latest[h] = set
	}

	var results []DominanceResult
	for _, h := range hosts {
		hSet, ok := latest[h]
		if !ok || len(hSet) < minPorts {
			continue
		}
		subsetCount := 0
		peers := 0
		for _, other := range hosts {
			if other == h {
				continue
			}
			oSet, ok2 := latest[other]
			if !ok2 {
				continue
			}
			peers++
			if isSubset(oSet, hSet) {
				subsetCount++
			}
		}
		score := 0.0
		if peers > 0 {
			score = float64(subsetCount) / float64(peers)
		}
		results = append(results, DominanceResult{
			Host:  h,
			Score: score,
			Peers: peers,
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

// isSubset returns true when every key in sub is present in super.
func isSubset(sub, super map[int]struct{}) bool {
	for k := range sub {
		if _, ok := super[k]; !ok {
			return false
		}
	}
	return true
}
