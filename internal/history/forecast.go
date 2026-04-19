package history

import (
	"sort"
	"time"
)

// ForecastEntry holds a predicted port open probability for a host.
type ForecastEntry struct {
	Host        string
	Port        int
	SeenCount   int
	TotalScans  int
	Probability float64 // 0.0 - 1.0
}

// Forecast predicts which ports are likely to be open based on historical
// frequency. Only ports seen in at least minCount scans are included.
func Forecast(store Store, host string, minCount int) []ForecastEntry {
	entries := store.ForHost(host)
	if len(entries) == 0 {
		return nil
	}

	totalScans := len(entries)
	portSeen := make(map[int]int)

	for _, e := range entries {
		for _, p := range e.Ports {
			portSeen[p]++
		}
	}

	var results []ForecastEntry
	for port, count := range portSeen {
		if count < minCount {
			continue
		}
		results = append(results, ForecastEntry{
			Host:        host,
			Port:        port,
			SeenCount:   count,
			TotalScans:  totalScans,
			Probability: float64(count) / float64(totalScans),
		})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Probability != results[j].Probability {
			return results[i].Probability > results[j].Probability
		}
		return results[i].Port < results[j].Port
	})

	return results
}

// ForecastSince limits the forecast to entries after the given time.
func ForecastSince(store Store, host string, since time.Time, minCount int) []ForecastEntry {
	filtered := Filter(store, FilterOptions{Host: host, Since: since})
	return Forecast(filtered, host, minCount)
}
