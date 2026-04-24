package history

import (
	"sort"
	"time"
)

// CohortEntry represents a group of hosts that first appeared with a common port
// set within the same time bucket.
type CohortEntry struct {
	CohortID  string   // e.g. "2024-01"
	Hosts     []string
	Ports     []int
	FirstSeen time.Time
}

// Cohort groups hosts by the time period (month) in which they were first
// observed, then intersects their open port sets to find shared ports.
func Cohort(store Store, minHosts int) []CohortEntry {
	allHosts := store.Hosts()
	type bucket struct {
		firstSeen time.Time
		hosts     []string
		portSets  []map[int]bool
	}
	buckets := map[string]*bucket{}

	for _, host := range allHosts {
		entries := store.ForHost(host)
		if len(entries) == 0 {
			continue
		}
		first := entries[0]
		key := first.Time.Format("2006-01")
		b, ok := buckets[key]
		if !ok {
			b = &bucket{firstSeen: first.Time}
			buckets[key] = b
		}
		if first.Time.Before(b.firstSeen) {
			b.firstSeen = first.Time
		}
		ps := make(map[int]bool, len(first.Ports))
		for _, p := range first.Ports {
			ps[p] = true
		}
		b.hosts = append(b.hosts, host)
		b.portSets = append(b.portSets, ps)
	}

	var result []CohortEntry
	for key, b := range buckets {
		if len(b.hosts) < minHosts {
			continue
		}
		shared := intersectSets(b.portSets)
		ports := make([]int, 0, len(shared))
		for p := range shared {
			ports = append(ports, p)
		}
		sort.Ints(ports)
		sort.Strings(b.hosts)
		result = append(result, CohortEntry{
			CohortID:  key,
			Hosts:     b.hosts,
			Ports:     ports,
			FirstSeen: b.firstSeen,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CohortID < result[j].CohortID
	})
	return result
}

func intersectSets(sets []map[int]bool) map[int]bool {
	if len(sets) == 0 {
		return map[int]bool{}
	}
	result := make(map[int]bool)
	for p := range sets[0] {
		result[p] = true
	}
	for _, s := range sets[1:] {
		for p := range result {
			if !s[p] {
				delete(result, p)
			}
		}
	}
	return result
}
