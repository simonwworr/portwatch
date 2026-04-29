package history

import "math"

// FluxResult holds the port flux score for a single host.
// Flux measures how rapidly the set of open ports changes between
// consecutive scans, expressed as the average symmetric-difference
// size normalised by the average union size (a Jaccard-distance mean).
type FluxResult struct {
	Host     string  `json:"host"`
	Flux     float64 `json:"flux"`      // 0.0 (stable) – 1.0 (fully volatile)
	AvgDelta float64 `json:"avg_delta"` // mean number of ports that changed per scan
	Scans    int     `json:"scans"`
}

// Flux computes per-host port-flux scores from the provided Store.
// minScans is the minimum number of scan entries required for a host
// to be included in the results.
func Flux(st Store, minScans int) []FluxResult {
	hosts := st.Hosts()
	var results []FluxResult

	for _, host := range hosts {
		entries := st.ForHost(host)
		if len(entries) < minScans {
			continue
		}

		var totalJaccard float64
		var totalDelta float64
		pairs := 0

		for i := 1; i < len(entries); i++ {
			prev := toIntSet(entries[i-1].Ports)
			curr := toIntSet(entries[i].Ports)

			union := unionSets(prev, curr)
			symDiff := symDiffSize(prev, curr)

			if len(union) == 0 {
				continue
			}

			totalJaccard += float64(symDiff) / float64(len(union))
			totalDelta += float64(symDiff)
			pairs++
		}

		if pairs == 0 {
			continue
		}

		flux := totalJaccard / float64(pairs)
		avgDelta := totalDelta / float64(pairs)

		results = append(results, FluxResult{
			Host:     host,
			Flux:     math.Round(flux*1000) / 1000,
			AvgDelta: math.Round(avgDelta*100) / 100,
			Scans:    len(entries),
		})
	}

	// sort descending by Flux
	for i := 1; i < len(results); i++ {
		for j := i; j > 0 && results[j].Flux > results[j-1].Flux; j-- {
			results[j], results[j-1] = results[j-1], results[j]
		}
	}

	return results
}
