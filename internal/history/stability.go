package history

import "sort"

// StabilityResult holds the stability score for a single host.
type StabilityResult struct {
	Host        string
	Scans       int
	Changes     int
	Stability   float64 // 1.0 = perfectly stable, 0.0 = changes every scan
	StablePorts []int
}

// Stability computes a stability score per host based on how often its port
// set changes between consecutive scans. Hosts with fewer than minScans
// entries are excluded.
func Stability(store Store, minScans int) []StabilityResult {
	hosts := store.Hosts()
	results := make([]StabilityResult, 0, len(hosts))

	for _, host := range hosts {
		entries := store.ForHost(host)
		if len(entries) < minScans {
			continue
		}

		changes := 0
		for i := 1; i < len(entries); i++ {
			prev := toIntSet(entries[i-1].Ports)
			curr := toIntSet(entries[i].Ports)
			if !intSetsEqual(prev, curr) {
				changes++
			}
		}

		scans := len(entries)
		possible := scans - 1
		var score float64
		if possible > 0 {
			score = 1.0 - float64(changes)/float64(possible)
		} else {
			score = 1.0
		}

		// stable ports: present in every scan
		stable := stablePortsFor(entries)

		results = append(results, StabilityResult{
			Host:        host,
			Scans:       scans,
			Changes:     changes,
			Stability:   score,
			StablePorts: stable,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Stability != results[j].Stability {
			return results[i].Stability > results[j].Stability
		}
		return results[i].Host < results[j].Host
	})
	return results
}

func toIntSet(ports []int) map[int]struct{} {
	s := make(map[int]struct{}, len(ports))
	for _, p := range ports {
		s[p] = struct{}{}
	}
	return s
}

func intSetsEqual(a, b map[int]struct{}) bool {
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

func stablePortsFor(entries []Entry) []int {
	if len(entries) == 0 {
		return nil
	}
	counts := map[int]int{}
	for _, e := range entries {
		for _, p := range e.Ports {
			counts[p]++
		}
	}
	total := len(entries)
	var stable []int
	for p, c := range counts {
		if c == total {
			stable = append(stable, p)
		}
	}
	sort.Ints(stable)
	return stable
}
