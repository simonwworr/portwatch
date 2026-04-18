package history

import "sort"

// RollupEntry summarizes port activity for a host over a time window.
type RollupEntry struct {
	Host      string
	ScanCount int
	UniquePorts []int
	MaxOpen   int
	MinOpen   int
}

// Rollup aggregates history entries per host into summary statistics.
func Rollup(store Store) []RollupEntry {
	hosts := make(map[string][][]int)
	for host, entries := range store {
		for _, e := range entries {
			hosts[host] = append(hosts[host], e.Ports)
		}
	}

	var result []RollupEntry
	for host, scans := range hosts {
		if len(scans) == 0 {
			continue
		}
		unique := make(map[int]struct{})
		maxOpen := -1
		minOpen := int(^uint(0) >> 1)
		for _, ports := range scans {
			for _, p := range ports {
				unique[p] = struct{}{}
			}
			if len(ports) > maxOpen {
				maxOpen = len(ports)
			}
			if len(ports) < minOpen {
				minOpen = len(ports)
			}
		}
		var up []int
		for p := range unique {
			up = append(up, p)
		}
		sort.Ints(up)
		result = append(result, RollupEntry{
			Host:        host,
			ScanCount:   len(scans),
			UniquePorts: up,
			MaxOpen:     maxOpen,
			MinOpen:     minOpen,
		})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Host < result[j].Host })
	return result
}
