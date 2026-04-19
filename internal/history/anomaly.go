package history

import "sort"

// AnomalyResult holds ports that appear rarely across scans for a host.
type AnomalyResult struct {
	Host        string
	RarePorts   []int
	ThresholdPct float64
}

// DetectAnomalies finds ports that appeared in fewer than thresholdPct percent
// of scans for each host in the store.
func DetectAnomalies(s *Store, thresholdPct float64) []AnomalyResult {
	var results []AnomalyResult

	for _, host := range s.Hosts() {
		entries := s.ForHost(host)
		if len(entries) == 0 {
			continue
		}

		freq := map[int]int{}
		for _, e := range entries {
			for _, p := range e.Ports {
				freq[p]++
			}
		}

		total := len(entries)
		var rare []int
		for port, count := range freq {
			pct := float64(count) / float64(total) * 100.0
			if pct < thresholdPct {
				rare = append(rare, port)
			}
		}
		sort.Ints(rare)

		if len(rare) > 0 {
			results = append(results, AnomalyResult{
				Host:        host,
				RarePorts:   rare,
				ThresholdPct: thresholdPct,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Host < results[j].Host
	})
	return results
}
