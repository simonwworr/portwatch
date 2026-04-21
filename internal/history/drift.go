package history

import "sort"

// DriftResult describes how a host's open ports have shifted over time
// relative to its earliest recorded snapshot.
type DriftResult struct {
	Host      string
	Added     []int // ports present now but not in baseline
	Removed   []int // ports in baseline but no longer present
	Stable    []int // ports present in both
	DriftScore float64 // (added+removed) / max(baseline, current)
}

// Drift compares the first and last snapshots for every host in the store
// and returns a DriftResult per host, ordered by DriftScore descending.
func Drift(store Store) []DriftResult {
	hosts := store.Hosts()
	var results []DriftResult

	for _, host := range hosts {
		entries := store.ForHost(host)
		if len(entries) < 2 {
			continue
		}

		baseline := toSet(entries[0].Ports)
		latest := toSet(entries[len(entries)-1].Ports)

		var added, removed, stable []int

		for p := range latest {
			if baseline[p] {
				stable = append(stable, p)
			} else {
				added = append(added, p)
			}
		}
		for p := range baseline {
			if !latest[p] {
				removed = append(removed, p)
			}
		}

		sort.Ints(added)
		sort.Ints(removed)
		sort.Ints(stable)

		denom := len(baseline)
		if len(latest) > denom {
			denom = len(latest)
		}
		var score float64
		if denom > 0 {
			score = float64(len(added)+len(removed)) / float64(denom)
		}

		results = append(results, DriftResult{
			Host:       host,
			Added:      added,
			Removed:    removed,
			Stable:     stable,
			DriftScore: score,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].DriftScore > results[j].DriftScore
	})
	return results
}
