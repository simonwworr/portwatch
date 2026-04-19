package history

import "sort"

// PortPattern describes how often a port appears across scans for a host.
type PortPattern struct {
	Host      string
	Port      int
	SeenCount int
	TotalScans int
	Frequency float64 // SeenCount / TotalScans
}

// Pattern analyses the scan history and returns port appearance frequencies
// for each host. Only hosts with at least minScans entries are included.
func Pattern(store Store, minScans int) []PortPattern {
	hosts := store.Hosts()
	var results []PortPattern

	for _, host := range hosts {
		entries := store.ForHost(host)
		if len(entries) < minScans {
			continue
		}
		total := len(entries)
		counts := map[int]int{}
		for _, e := range entries {
			for _, p := range e.Ports {
				counts[p]++
			}
		}
		for port, seen := range counts {
			results = append(results, PortPattern{
				Host:       host,
				Port:       port,
				SeenCount:  seen,
				TotalScans: total,
				Frequency:  float64(seen) / float64(total),
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Host != results[j].Host {
			return results[i].Host < results[j].Host
		}
		return results[i].Frequency > results[j].Frequency
	})
	return results
}
