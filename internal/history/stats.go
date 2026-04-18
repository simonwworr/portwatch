package history

import "sort"

// PortStats holds frequency statistics for a single port.
type PortStats struct {
	Port      int
	SeenCount int
}

// HostStats holds aggregated stats for a host.
type HostStats struct {
	Host      string
	ScanCount int
	TopPorts  []PortStats
}

// Stats computes scan frequency and top open ports for each host in the store.
// topN controls how many ports are returned (0 = all).
func Stats(s *Store, topN int) []HostStats {
	type portFreq map[int]int
	freqs := make(map[string]portFreq)
	scans := make(map[string]int)

	for host, entries := range s.entries {
		scans[host] = len(entries)
		if freqs[host] == nil {
			freqs[host] = make(portFreq)
		}
		for _, e := range entries {
			for _, p := range e.Ports {
				freqs[host][p]++
			}
		}
	}

	var result []HostStats
	for host, pf := range freqs {
		var ps []PortStats
		for port, count := range pf {
			ps = append(ps, PortStats{Port: port, SeenCount: count})
		}
		sort.Slice(ps, func(i, j int) bool {
			if ps[i].SeenCount != ps[j].SeenCount {
				return ps[i].SeenCount > ps[j].SeenCount
			}
			return ps[i].Port < ps[j].Port
		})
		if topN > 0 && len(ps) > topN {
			ps = ps[:topN]
		}
		result = append(result, HostStats{
			Host:      host,
			ScanCount: scans[host],
			TopPorts:  ps,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Host < result[j].Host
	})
	return result
}
