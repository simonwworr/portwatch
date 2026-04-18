package history

import "sort"

// PortDelta represents the change in open ports for a host between two points in time.
type PortDelta struct {
	Host    string
	Opened  []int
	Closed  []int
	Stable  []int
}

// Compare returns the port deltas between the two most recent scan entries for
// each host found in store. Hosts with fewer than two entries are skipped.
func Compare(store *Store) []PortDelta {
	hosts := store.Hosts()
	result := make([]PortDelta, 0, len(hosts))

	for _, host := range hosts {
		entries := store.ForHost(host)
		if len(entries) < 2 {
			continue
		}
		prev := toPortSet(entries[len(entries)-2].Ports)
		curr := toPortSet(entries[len(entries)-1].Ports)

		delta := PortDelta{Host: host}
		for p := range curr {
			if prev[p] {
				delta.Stable = append(delta.Stable, p)
			} else {
				delta.Opened = append(delta.Opened, p)
			}
		}
		for p := range prev {
			if !curr[p] {
				delta.Closed = append(delta.Closed, p)
			}
		}
		sort.Ints(delta.Opened)
		sort.Ints(delta.Closed)
		sort.Ints(delta.Stable)
		result = append(result, delta)
	}
	return result
}

func toPortSet(ports []int) map[int]bool {
	s := make(map[int]bool, len(ports))
	for _, p := range ports {
		s[p] = true
	}
	return s
}
