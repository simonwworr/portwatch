package history

import (
	"math"
	"sort"
)

// MomentumResult holds the momentum score for a host.
// A positive score indicates ports are trending open; negative means closing.
type MomentumResult struct {
	Host      string
	Score     float64 // weighted rate of change: positive = opening, negative = closing
	NetChange int     // total net port change over window
	Scans     int
}

// Momentum computes a weighted rate-of-change score for each host.
// Recent changes are weighted more heavily using exponential decay (lambda).
// minScans is the minimum number of entries required to compute a score.
func Momentum(store Store, minScans int, lambda float64) []MomentumResult {
	if lambda <= 0 {
		lambda = 0.5
	}

	hosts := allHosts(store)
	var results []MomentumResult

	for _, host := range hosts {
		entries, _ := store.ForHost(host)
		if len(entries) < minScans {
			continue
		}

		// Sort ascending by time
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Time.Before(entries[j].Time)
		})

		n := len(entries)
		var weightedSum float64
		var totalWeight float64
		netChange := 0

		for i := 1; i < n; i++ {
			delta := len(entries[i].Ports) - len(entries[i-1].Ports)
			// Weight by recency: most recent diff gets weight 1, older get exponential decay
			age := float64(n - 1 - i)
			weight := math.Exp(-lambda * age)
			weightedSum += float64(delta) * weight
			totalWeight += weight
			netChange += delta
		}

		score := 0.0
		if totalWeight > 0 {
			score = weightedSum / totalWeight
		}

		results = append(results, MomentumResult{
			Host:      host,
			Score:     math.Round(score*1000) / 1000,
			NetChange: netChange,
			Scans:     n,
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
