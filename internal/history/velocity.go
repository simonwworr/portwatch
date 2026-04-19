package history

import (
	"sort"
	"time"
)

// VelocityEntry holds the rate of port change for a host over a window.
type VelocityEntry struct {
	Host        string
	OpenRate    float64 // ports opened per day
	CloseRate   float64 // ports closed per day
	NetRate     float64 // net change per day
	WindowDays  int
	ScanCount   int
}

// Velocity calculates the rate of port change per host over the given window.
func Velocity(store *Store, window time.Duration) []VelocityEntry {
	since := time.Now().Add(-window)
	windowDays := window.Hours() / 24
	if windowDays < 1 {
		windowDays = 1
	}

	hosts := store.Hosts()
	results := make([]VelocityEntry, 0, len(hosts))

	for _, host := range hosts {
		entries := store.ForHost(host)
		var windowed []Entry
		for _, e := range entries {
			if !e.Time.Before(since) {
				windowed = append(windowed, e)
			}
		}
		if len(windowed) < 2 {
			continue
		}

		sort.Slice(windowed, func(i, j int) bool {
			return windowed[i].Time.Before(windowed[j].Time)
		})

		var totalOpened, totalClosed int
		for i := 1; i < len(windowed); i++ {
			prev := toSet(windowed[i-1].Ports)
			curr := toSet(windowed[i].Ports)
			for p := range curr {
				if !prev[p] {
					totalOpened++
				}
			}
			for p := range prev {
				if !curr[p] {
					totalClosed++
				}
			}
		}

		openRate := float64(totalOpened) / windowDays
		closeRate := float64(totalClosed) / windowDays
		results = append(results, VelocityEntry{
			Host:       host,
			OpenRate:   openRate,
			CloseRate:  closeRate,
			NetRate:    openRate - closeRate,
			WindowDays: int(windowDays),
			ScanCount:  len(windowed),
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].NetRate > results[j].NetRate
	})
	return results
}
