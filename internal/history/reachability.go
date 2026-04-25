package history

import "sort"

// ReachabilityResult holds reachability metrics for a single host.
type ReachabilityResult struct {
	Host        string  `json:"host"`
	TotalScans  int     `json:"total_scans"`
	ActiveScan  int     `json:"active_scans"`
	Uptime      float64 `json:"uptime"`  // fraction of scans with at least one open port
	AvgPorts    float64 `json:"avg_ports"`
	MaxPorts    int     `json:"max_ports"`
	MinPorts    int     `json:"min_ports"`
}

// Reachability computes per-host uptime and port availability metrics.
// A host is considered "active" in a scan if it has at least one open port.
// minScans skips hosts with fewer than that many scan entries.
func Reachability(store Store, minScans int) []ReachabilityResult {
	hosts := store.Hosts()
	results := make([]ReachabilityResult, 0, len(hosts))

	for _, host := range hosts {
		entries := store.ForHost(host)
		if len(entries) < minScans {
			continue
		}

		active := 0
		totalPorts := 0
		maxP := 0
		minP := -1

		for _, e := range entries {
			n := len(e.Ports)
			if n > 0 {
				active++
			}
			totalPorts += n
			if n > maxP {
				maxP = n
			}
			if minP < 0 || n < minP {
				minP = n
			}
		}

		if minP < 0 {
			minP = 0
		}

		uptime := float64(active) / float64(len(entries))
		avgPorts := float64(totalPorts) / float64(len(entries))

		results = append(results, ReachabilityResult{
			Host:       host,
			TotalScans: len(entries),
			ActiveScan: active,
			Uptime:     uptime,
			AvgPorts:   avgPorts,
			MaxPorts:   maxP,
			MinPorts:   minP,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Uptime != results[j].Uptime {
			return results[i].Uptime > results[j].Uptime
		}
		return results[i].Host < results[j].Host
	})

	return results
}
