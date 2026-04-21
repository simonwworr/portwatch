package history

import "sort"

// HotspotEntry represents a port that has appeared frequently across multiple hosts.
type HotspotEntry struct {
	Port      int      `json:"port"`
	HostCount int      `json:"host_count"`
	Hosts     []string `json:"hosts"`
	Frequency float64  `json:"frequency"` // fraction of scans where port was seen, averaged across hosts
}

// Hotspot identifies ports that are consistently open across many hosts.
// minHosts filters out ports seen on fewer than minHosts distinct hosts.
func Hotspot(store Store, minHosts int) []HotspotEntry {
	all := store.All()
	if len(all) == 0 {
		return nil
	}

	// per-host: count total scans and scans where each port appeared
	type hostStats struct {
		total int
		ports map[int]int
	}
	hostMap := map[string]*hostStats{}
	for _, e := range all {
		hs, ok := hostMap[e.Host]
		if !ok {
			hs = &hostStats{ports: map[int]int{}}
			hostMap[e.Host] = hs
		}
		hs.total++
		seen := map[int]bool{}
		for _, p := range e.Ports {
			if !seen[p] {
				hs.ports[p]++
				seen[p] = true
			}
		}
	}

	// aggregate per port across hosts
	type portAgg struct {
		hosts     []string
		freqSum   float64
	}
	portMap := map[int]*portAgg{}
	for host, hs := range hostMap {
		if hs.total == 0 {
			continue
		}
		for port, count := range hs.ports {
			pa, ok := portMap[port]
			if !ok {
				pa = &portAgg{}
				portMap[port] = pa
			}
			pa.hosts = append(pa.hosts, host)
			pa.freqSum += float64(count) / float64(hs.total)
		}
	}

	var results []HotspotEntry
	for port, pa := range portMap {
		if len(pa.hosts) < minHosts {
			continue
		}
		sort.Strings(pa.hosts)
		results = append(results, HotspotEntry{
			Port:      port,
			HostCount: len(pa.hosts),
			Hosts:     pa.hosts,
			Frequency: pa.freqSum / float64(len(pa.hosts)),
		})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].HostCount != results[j].HostCount {
			return results[i].HostCount > results[j].HostCount
		}
		return results[i].Port < results[j].Port
	})
	return results
}
