package history

import "time"

// HostSummary holds aggregated statistics for a single host.
type HostSummary struct {
	Host        string
	TotalScans  int
	FirstSeen   time.Time
	LastSeen    time.Time
	UniquePorts []int
	MaxOpen     int
	MinOpen     int
}

// Summarize computes a HostSummary for the given host from the store.
func Summarize(s *Store, host string) (HostSummary, bool) {
	entries := s.ForHost(host)
	if len(entries) == 0 {
		return HostSummary{}, false
	}

	portSet := map[int]struct{}{}
	minOpen := -1
	maxOpen := 0

	sum := HostSummary{
		Host:      host,
		FirstSeen: entries[0].ScannedAt,
		LastSeen:  entries[len(entries)-1].ScannedAt,
	}

	for _, e := range entries {
		sum.TotalScans++
		count := len(e.OpenPorts)
		if count > maxOpen {
			maxOpen = count
		}
		if minOpen == -1 || count < minOpen {
			minOpen = count
		}
		for _, p := range e.OpenPorts {
			portSet[p] = struct{}{}
		}
	}

	sum.MaxOpen = maxOpen
	if minOpen == -1 {
		minOpen = 0
	}
	sum.MinOpen = minOpen

	unique := make([]int, 0, len(portSet))
	for p := range portSet {
		unique = append(unique, p)
	}
	sortInts(unique)
	sum.UniquePorts = unique

	return sum, true
}

func sortInts(s []int) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
