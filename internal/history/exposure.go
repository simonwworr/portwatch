package history

import "sort"

// ExposureResult holds the cumulative open-time stats for a single port on a host.
type ExposureResult struct {
	Host       string
	Port       int
	SeenCount  int     // number of scans where port was open
	TotalScans int     // total scans recorded for this host
	Exposure   float64 // fraction of scans where port was open (0.0–1.0)
}

// Exposure calculates how frequently each port was observed open across all
// recorded scans for every host in the store. Only hosts with at least
// minScans entries are included.
func Exposure(store Store, minScans int) []ExposureResult {
	hosts := store.Hosts()
	var results []ExposureResult

	for _, host := range hosts {
		entries := store.ForHost(host)
		if len(entries) < minScans {
			continue
		}
		total := len(entries)

		// count per-port open occurrences
		counts := make(map[int]int)
		for _, e := range entries {
			for _, p := range e.Ports {
				counts[p]++
			}
		}

		for port, seen := range counts {
			results = append(results, ExposureResult{
				Host:       host,
				Port:       port,
				SeenCount:  seen,
				TotalScans: total,
				Exposure:   float64(seen) / float64(total),
			})
		}
	}

	// stable sort: descending exposure, then host, then port
	sort.Slice(results, func(i, j int) bool {
		if results[i].Exposure != results[j].Exposure {
			return results[i].Exposure > results[j].Exposure
		}
		if results[i].Host != results[j].Host {
			return results[i].Host < results[j].Host
		}
		return results[i].Port < results[j].Port
	})

	return results
}
