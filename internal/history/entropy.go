package history

import (
	"math"
	"sort"
)

// EntropyResult holds the Shannon entropy score for a host's port distribution.
type EntropyResult struct {
	Host    string
	Entropy float64
	Scans   int
	Unique  int
}

// Entropy computes the Shannon entropy of open-port distributions per host.
// Higher entropy indicates more varied/unpredictable port behaviour.
func Entropy(store Store, minScans int) []EntropyResult {
	hosts := store.Hosts()
	var results []EntropyResult

	for _, host := range hosts {
		entries := store.ForHost(host)
		if len(entries) < minScans {
			continue
		}

		// Count how often each port appears across all scans.
		freq := make(map[int]int)
		total := 0
		for _, e := range entries {
			for _, p := range e.Ports {
				freq[p]++
				total++
			}
		}

		if total == 0 {
			results = append(results, EntropyResult{Host: host, Scans: len(entries)})
			continue
		}

		// Shannon entropy: H = -sum(p * log2(p))
		var h float64
		for _, count := range freq {
			p := float64(count) / float64(total)
			h -= p * math.Log2(p)
		}

		results = append(results, EntropyResult{
			Host:    host,
			Entropy: math.Round(h*1000) / 1000,
			Scans:   len(entries),
			Unique:  len(freq),
		})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Entropy != results[j].Entropy {
			return results[i].Entropy > results[j].Entropy
		}
		return results[i].Host < results[j].Host
	})

	return results
}
