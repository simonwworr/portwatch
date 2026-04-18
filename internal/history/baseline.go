package history

import (
	"fmt"
	"sort"
	"time"
)

// Baseline represents the most recently known open ports for each host.
type Baseline map[string][]int

// BuildBaseline constructs a baseline from the latest scan entry per host.
func BuildBaseline(store *Store) Baseline {
	baseline := make(Baseline)
	for _, host := range store.Hosts() {
		entries := store.ForHost(host)
		if len(entries) == 0 {
			continue
		}
		// entries are ordered oldest-first; take the last
		latest := entries[len(entries)-1]
		ports := make([]int, len(latest.Ports))
		copy(ports, latest.Ports)
		sort.Ints(ports)
		baseline[host] = ports
	}
	return baseline
}

// Deviation describes ports that have opened or closed relative to a baseline.
type Deviation struct {
	Host    string
	Opened  []int
	Closed  []int
	ScannedAt time.Time
}

// DetectDeviations compares current scan results against a baseline.
func DetectDeviations(baseline Baseline, current map[string][]int, scannedAt time.Time) []Deviation {
	var result []Deviation
	for host, currentPorts := range current {
		base := baseline[host]
		bSet := toPortSet(base)
		cSet := toPortSet(currentPorts)

		var opened, closed []int
		for p := range cSet {
			if !bSet[p] {
				opened = append(opened, p)
			}
		}
		for p := range bSet {
			if !cSet[p] {
				closed = append(closed, p)
			}
		}
		if len(opened) > 0 || len(closed) > 0 {
			sort.Ints(opened)
			sort.Ints(closed)
			result = append(result, Deviation{
				Host:      host,
				Opened:    opened,
				Closed:    closed,
				ScannedAt: scannedAt,
			})
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Host < result[j].Host })
	return result
}

func (d Deviation) String() string {
	return fmt.Sprintf("host=%s opened=%v closed=%v at=%s", d.Host, d.Opened, d.Closed, d.ScannedAt.Format(time.RFC3339))
}
