package history

import "sort"

// SequenceEntry describes a consecutive run of scans where a port remained open.
type SequenceEntry struct {
	Host  string
	Port  int
	Start string // RFC3339 timestamp of first seen
	End   string // RFC3339 timestamp of last seen
	Runs  int    // number of consecutive scans
}

// Sequence finds ports that have appeared in consecutive scans for each host,
// returning only those with at least minRuns consecutive appearances.
func Sequence(store Store, minRuns int) []SequenceEntry {
	if minRuns < 2 {
		minRuns = 2
	}

	hosts := store.Hosts()
	var results []SequenceEntry

	for _, host := range hosts {
		entries := store.ForHost(host)
		if len(entries) < minRuns {
			continue
		}
		// sort by time ascending
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Time.Before(entries[j].Time)
		})

		// track run length per port
		type run struct {
			start string
			last  string
			count int
		}
		runs := map[int]*run{}

		for _, e := range entries {
			portSet := map[int]bool{}
			for _, p := range e.Ports {
				portSet[p] = true
			}
			ts := e.Time.Format("2006-01-02T15:04:05Z07:00")
			// update runs
			for p, r := range runs {
				if portSet[p] {
					r.last = ts
					r.count++
				} else {
					if r.count >= minRuns {
						results = append(results, SequenceEntry{Host: host, Port: p, Start: r.start, End: r.last, Runs: r.count})
					}
					delete(runs, p)
				}
			}
			for p := range portSet {
				if _, ok := runs[p]; !ok {
					runs[p] = &run{start: ts, last: ts, count: 1}
				}
			}
		}
		for p, r := range runs {
			if r.count >= minRuns {
				results = append(results, SequenceEntry{Host: host, Port: p, Start: r.start, End: r.last, Runs: r.count})
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
