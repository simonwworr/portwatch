package history

import "math"

// VolatilityResult holds the volatility score for a single host.
type VolatilityResult struct {
	Host       string  `json:"host"`
	Scans      int     `json:"scans"`
	Changes    int     `json:"changes"`
	AvgDelta   float64 `json:"avg_delta"`
	Volatility float64 `json:"volatility"`
}

// Volatility measures how much the open port set changes between consecutive
// scans for each host. The volatility score is the mean symmetric difference
// size (normalised by the union size — Jaccard distance) across all scan pairs.
func Volatility(store Store, minScans int) []VolatilityResult {
	hosts := store.Hosts()
	var results []VolatilityResult

	for _, host := range hosts {
		entries := store.ForHost(host)
		if len(entries) < minScans {
			continue
		}

		var totalDist float64
		changePairs := 0

		for i := 1; i < len(entries); i++ {
			prev := toIntSet(entries[i-1].Ports)
			curr := toIntSet(entries[i].Ports)

			unionSize := len(unionSets(prev, curr))
			if unionSize == 0 {
				continue
			}

			symDiff := symDiffSize(prev, curr)
			dist := float64(symDiff) / float64(unionSize)
			totalDist += dist
			if symDiff > 0 {
				changePairs++
			}
		}

		pairs := len(entries) - 1
		if pairs == 0 {
			continue
		}

		avgDelta := totalDist / float64(pairs)
		volScore := math.Round(avgDelta*1000) / 1000

		results = append(results, VolatilityResult{
			Host:       host,
			Scans:      len(entries),
			Changes:    changePairs,
			AvgDelta:   math.Round(avgDelta*1000) / 1000,
			Volatility: volScore,
		})
	}

	return results
}

func unionSets(a, b map[int]struct{}) map[int]struct{} {
	u := make(map[int]struct{}, len(a)+len(b))
	for k := range a {
		u[k] = struct{}{}
	}
	for k := range b {
		u[k] = struct{}{}
	}
	return u
}

func symDiffSize(a, b map[int]struct{}) int {
	count := 0
	for k := range a {
		if _, ok := b[k]; !ok {
			count++
		}
	}
	for k := range b {
		if _, ok := a[k]; !ok {
			count++
		}
	}
	return count
}
