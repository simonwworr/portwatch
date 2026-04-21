package history

import "sort"

// OverlapResult describes the port overlap between two hosts.
type OverlapResult struct {
	HostA        string
	HostB        string
	SharedPorts  []int
	OnlyInA      []int
	OnlyInB      []int
	JaccardScore float64 // |A ∩ B| / |A ∪ B|
}

// Overlap computes the port overlap between every pair of hosts using their
// latest snapshot. Results are sorted descending by JaccardScore.
func Overlap(store Store) []OverlapResult {
	snaps := AllSnapshots(store)
	if len(snaps) < 2 {
		return nil
	}

	// Build port sets keyed by host.
	sets := make(map[string]map[int]struct{}, len(snaps))
	hosts := make([]string, 0, len(snaps))
	for _, s := range snaps {
		set := make(map[int]struct{}, len(s.Ports))
		for _, p := range s.Ports {
			set[p] = struct{}{}
		}
		sets[s.Host] = set
		hosts = append(hosts, s.Host)
	}
	sort.Strings(hosts)

	var results []OverlapResult
	for i := 0; i < len(hosts); i++ {
		for j := i + 1; j < len(hosts); j++ {
			a, b := hosts[i], hosts[j]
			setA, setB := sets[a], sets[b]

			var shared, onlyA, onlyB []int
			for p := range setA {
				if _, ok := setB[p]; ok {
					shared = append(shared, p)
				} else {
					onlyA = append(onlyA, p)
				}
			}
			for p := range setB {
				if _, ok := setA[p]; !ok {
					onlyB = append(onlyB, p)
				}
			}

			unionSize := len(shared) + len(onlyA) + len(onlyB)
			var jaccard float64
			if unionSize > 0 {
				jaccard = float64(len(shared)) / float64(unionSize)
			}

			sort.Ints(shared)
			sort.Ints(onlyA)
			sort.Ints(onlyB)

			results = append(results, OverlapResult{
				HostA:        a,
				HostB:        b,
				SharedPorts:  shared,
				OnlyInA:      onlyA,
				OnlyInB:      onlyB,
				JaccardScore: jaccard,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].JaccardScore > results[j].JaccardScore
	})
	return results
}
