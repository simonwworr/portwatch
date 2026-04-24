package history

import (
	"sort"
	"time"
)

// ShadowPort represents a port that appeared briefly and then disappeared.
type ShadowPort struct {
	Host     string
	Port     int
	FirstSeen time.Time
	LastSeen  time.Time
	Appearances int
}

// Shadow detects ports that appear for fewer than minScans scans across all
// entries for each host, flagging them as transient or "shadow" ports.
// A host must have at least minScans total scan entries to be considered;
// hosts with fewer entries are skipped to avoid false positives on new hosts.
// A port is flagged as a shadow port if it appeared in at most maxAppearances
// of the scans recorded for that host.
func Shadow(store Store, minScans int, maxAppearances int) []ShadowPort {
	hosts := store.Hosts()
	var results []ShadowPort

	for _, host := range hosts {
		entries := store.ForHost(host)
		if len(entries) < minScans {
			continue
		}

		type portMeta struct {
			first time.Time
			last  time.Time
			count int
		}
		portMap := map[int]*portMeta{}

		for _, e := range entries {
			for _, p := range e.Ports {
				if m, ok := portMap[p]; ok {
					m.count++
					if e.Time.After(m.last) {
						m.last = e.Time
					}
				} else {
					portMap[p] = &portMeta{first: e.Time, last: e.Time, count: 1}
				}
			}
		}

		for port, m := range portMap {
			if m.count <= maxAppearances {
				results = append(results, ShadowPort{
					Host:        host,
					Port:        port,
					FirstSeen:   m.first,
					LastSeen:    m.last,
					Appearances: m.count,
				})
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Host != results[j].Host {
			return results[i].Host < results[j].Host
		}
		return results[i].Port < results[j].Port
	})
	return results
}

// Duration returns the time elapsed between the first and last time this
// shadow port was observed. Returns zero for ports seen only once.
func (s ShadowPort) Duration() time.Duration {
	return s.LastSeen.Sub(s.FirstSeen)
}
