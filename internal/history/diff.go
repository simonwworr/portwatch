package history

import "time"

// PortDiff represents the port changes between two consecutive scans for a host.
type PortDiff struct {
	Host    string
	At      time.Time
	Opened  []int
	Closed  []int
}

// DiffHistory computes port diffs for a host across all consecutive scan pairs.
func DiffHistory(store *Store, host string) []PortDiff {
	entries := store.ForHost(host)
	if len(entries) < 2 {
		return nil
	}

	var diffs []PortDiff
	for i := 1; i < len(entries); i++ {
		prev := toPortSetSlice(entries[i-1].Ports)
		curr := toPortSetSlice(entries[i].Ports)

		opened := setDiff(curr, prev)
		closed := setDiff(prev, curr)

		if len(opened) > 0 || len(closed) > 0 {
			diffs = append(diffs, PortDiff{
				Host:   host,
				At:     entries[i].Time,
				Opened: opened,
				Closed: closed,
			})
		}
	}
	return diffs
}

func toPortSetSlice(ports []int) map[int]struct{} {
	m := make(map[int]struct{}, len(ports))
	for _, p := range ports {
		m[p] = struct{}{}
	}
	return m
}

func setDiff(a, b map[int]struct{}) []int {
	var result []int
	for k := range a {
		if _, ok := b[k]; !ok {
			result = append(result, k)
		}
	}
	sortInts(result)
	return result
}
