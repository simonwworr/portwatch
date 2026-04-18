package history

import (
	"time"
)

// Snapshot represents the port state for a host at a point in time.
type Snapshot struct {
	Host      string    `json:"host"`
	Ports     []int     `json:"ports"`
	Timestamp time.Time `json:"timestamp"`
}

// LatestSnapshot returns the most recent entry for the given host.
func LatestSnapshot(store *Store, host string) (Snapshot, bool) {
	entries := store.ForHost(host)
	if len(entries) == 0 {
		return Snapshot{}, false
	}
	last := entries[len(entries)-1]
	return Snapshot{
		Host:      host,
		Ports:     last.Ports,
		Timestamp: last.Timestamp,
	}, true
}

// AllSnapshots returns the latest snapshot for every host in the store.
func AllSnapshots(store *Store) []Snapshot {
	seen := map[string]bool{}
	var out []Snapshot
	for _, host := range store.Hosts() {
		if seen[host] {
			continue
		}
		seen[host] = true
		if s, ok := LatestSnapshot(store, host); ok {
			out = append(out, s)
		}
	}
	return out
}
