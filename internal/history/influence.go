package history

import "sort"

// InfluenceResult describes how much a single host's port profile
// affects the rest of the fleet — measured as the average Jaccard
// similarity between that host and every other host that shares at
// least one port with it.
type InfluenceResult struct {
	Host        string  `json:"host"`
	Peers       int     `json:"peers"`        // number of hosts that share ≥1 port
	SharedPorts int     `json:"shared_ports"` // total distinct ports shared across all peers
	Score       float64 `json:"score"`        // mean Jaccard similarity to peers
}

// Influence computes, for each host in the store, how much its latest
// open-port fingerprint overlaps with every other host's fingerprint.
// Only hosts that have at least one scan entry are included.
// Results are sorted descending by Score.
func Influence(store Store) []InfluenceResult {
	// Collect the latest snapshot for every host.
	snaps := AllSnapshots(store)
	if len(snaps) == 0 {
		return nil
	}

	// Build a map of host → port-set for quick lookup.
	portSets := make(map[string]map[int]struct{}, len(snaps))
	for _, s := range snaps {
		set := make(map[int]struct{}, len(s.Ports))
		for _, p := range s.Ports {
			set[p] = struct{}{}
		}
		portSets[s.Host] = set
	}

	hosts := make([]string, 0, len(portSets))
	for h := range portSets {
		hosts = append(hosts, h)
	}
	sort.Strings(hosts)

	results := make([]InfluenceResult, 0, len(hosts))

	for _, h := range hosts {
		setA := portSets[h]
		if len(setA) == 0 {
			// Host with no open ports cannot influence anyone.
			results = append(results, InfluenceResult{Host: h})
			continue
		}

		var totalJaccard float64
		peers := 0
		sharedUnion := make(map[int]struct{})

		for _, other := range hosts {
			if other == h {
				continue
			}
			setB := portSets[other]

			// Intersection
			var inter int
			for p := range setA {
				if _, ok := setB[p]; ok {
					inter++
					sharedUnion[p] = struct{}{}
				}
			}
			if inter == 0 {
				continue // no overlap — not a peer
			}

			// Union
			union := len(setA) + len(setB) - inter
			if union > 0 {
				totalJaccard += float64(inter) / float64(union)
			}
			peers++
		}

		var score float64
		if peers > 0 {
			score = totalJaccard / float64(peers)
		}

		results = append(results, InfluenceResult{
			Host:        h,
			Peers:       peers,
			SharedPorts: len(sharedUnion),
			Score:       score,
		})
	}

	// Sort descending by Score, then by Host for determinism.
	sort.Slice(results, func(i, j int) bool {
		if results[i].Score != results[j].Score {
			return results[i].Score > results[j].Score
		}
		return results[i].Host < results[j].Host
	})

	return results
}
