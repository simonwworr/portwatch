package history

import "sort"

// PortTrend describes how a port's presence has changed over time.
type PortTrend struct {
	Port       int     `json:"port"`
	SeenCount  int     `json:"seen_count"`
	TotalScans int     `json:"total_scans"`
	Frequency  float64 `json:"frequency"`
	FirstSeen  string  `json:"first_seen"`
	LastSeen   string  `json:"last_seen"`
}

// HostTrend holds trend data for a single host.
type HostTrend struct {
	Host   string      `json:"host"`
	Trends []PortTrend `json:"trends"`
}

// Trend computes per-port frequency trends for a host across all history entries.
func Trend(store *Store, host string) HostTrend {
	entries := store.ForHost(host)
	total := len(entries)
	if total == 0 {
		return HostTrend{Host: host}
	}

	portSeen := map[int]int{}
	portFirst := map[int]string{}
	portLast := map[int]string{}

	for _, e := range entries {
		ts := e.ScannedAt.Format("2006-01-02T15:04:05Z")
		for _, p := range e.OpenPorts {
			portSeen[p]++
			if _, ok := portFirst[p]; !ok {
				portFirst[p] = ts
			}
			portLast[p] = ts
		}
	}

	var trends []PortTrend
	for port, count := range portSeen {
		trends = append(trends, PortTrend{
			Port:       port,
			SeenCount:  count,
			TotalScans: total,
			Frequency:  float64(count) / float64(total),
			FirstSeen:  portFirst[port],
			LastSeen:   portLast[port],
		})
	}

	sort.Slice(trends, func(i, j int) bool {
		return trends[i].Port < trends[j].Port
	})

	return HostTrend{Host: host, Trends: trends}
}
