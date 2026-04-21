package history

import "sort"

// ChurnResult holds port churn metrics for a single host.
type ChurnResult struct {
	Host       string
	Opened     int
	Closed     int
	Total      int
	ChurnRate  float64 // (opened + closed) / total scans
}

// Churn calculates how frequently ports open and close across scans for each host.
// minScans is the minimum number of scan entries required to compute churn.
func Churn(store Store, minScans int) []ChurnResult {
	hosts := store.Hosts()
	results := make([]ChurnResult, 0, len(hosts))

	for _, host := range hosts {
		entries := store.ForHost(host)
		if len(entries) < minScans {
			continue
		}

		var opened, closed int
		for i := 1; i < len(entries); i++ {
			prev := toPortSetSlice(entries[i-1].Ports)
			curr := toPortSetSlice(entries[i].Ports)
			opened += len(setDiff(curr, prev))
			closed += len(setDiff(prev, curr))
		}

		total := opened + closed
		rate := 0.0
		if len(entries) > 1 {
			rate = float64(total) / float64(len(entries)-1)
		}

		results = append(results, ChurnResult{
			Host:      host,
			Opened:    opened,
			Closed:    closed,
			Total:     total,
			ChurnRate: rate,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].ChurnRate > results[j].ChurnRate
	})

	return results
}
