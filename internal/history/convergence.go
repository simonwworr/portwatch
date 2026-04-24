package history

import "sort"

// ConvergenceResult describes how quickly a host's port set stabilises over time.
type ConvergenceResult struct {
	Host         string  `json:"host"`
	Scans        int     `json:"scans"`
	// StableAfter is the scan index (1-based) after which the port set no
	// longer changed until the end of the recorded history.
	StableAfter  int     `json:"stable_after"`
	// Stable is true when the port set did not change in the final window.
	Stable       bool    `json:"stable"`
	// ChangeRate is the fraction of consecutive scan pairs that differed.
	ChangeRate   float64 `json:"change_rate"`
}

// Convergence analyses each host's history and reports how quickly (or
// whether) its open-port set converges to a stable state.
func Convergence(store Store) []ConvergenceResult {
	hosts := store.Hosts()
	sort.Strings(hosts)

	var results []ConvergenceResult
	for _, host := range hosts {
		entries := store.ForHost(host)
		if len(entries) < 2 {
			continue
		}

		changes := 0
		lastChange := 0
		prev := toPortSet(entries[0].Ports)

		for i := 1; i < len(entries); i++ {
			curr := toPortSet(entries[i].Ports)
			if !setsEqual(prev, curr) {
				changes++
				lastChange = i
			}
			prev = curr
		}

		total := len(entries) - 1
		results = append(results, ConvergenceResult{
			Host:        host,
			Scans:       len(entries),
			StableAfter: lastChange + 1,
			Stable:      lastChange < total,
			ChangeRate:  float64(changes) / float64(total),
		})
	}
	return results
}

func setsEqual(a, b map[int]struct{}) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if _, ok := b[k]; !ok {
			return false
		}
	}
	return true
}
