package history

import "sort"

// CorrelationResult holds port co-occurrence data for a host.
type CorrelationResult struct {
	Host  string
	Pairs []PortPair
}

// PortPair represents two ports that appeared together and how often.
type PortPair struct {
	PortA int
	PortB int
	Count int
}

// Correlate finds ports that frequently appear together in the same scan
// across all entries for each host in the store.
func Correlate(store Store, minCount int) []CorrelationResult {
	byHost := map[string][][]int{}
	for _, e := range store.Entries {
		byHost[e.Host] = append(byHost[e.Host], e.Ports)
	}

	var results []CorrelationResult
	for host, scans := range byHost {
		counts := map[[2]int]int{}
		for _, ports := range scans {
			sorted := make([]int, len(ports))
			copy(sorted, ports)
			sort.Ints(sorted)
			for i := 0; i < len(sorted); i++ {
				for j := i + 1; j < len(sorted); j++ {
					key := [2]int{sorted[i], sorted[j]}
					counts[key]++
				}
			}
		}
		var pairs []PortPair
		for k, c := range counts {
			if c >= minCount {
				pairs = append(pairs, PortPair{PortA: k[0], PortB: k[1], Count: c})
			}
		}
		if len(pairs) == 0 {
			continue
		}
		sort.Slice(pairs, func(i, j int) bool {
			if pairs[i].Count != pairs[j].Count {
				return pairs[i].Count > pairs[j].Count
			}
			return pairs[i].PortA < pairs[j].PortA
		})
		results = append(results, CorrelationResult{Host: host, Pairs: pairs})
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Host < results[j].Host
	})
	return results
}
