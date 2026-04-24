package history

import "math"

// PressureResult holds the port pressure score for a single host.
// Pressure measures how many distinct ports are open on average, weighted
// by how frequently they appear, giving a sense of overall "load" or
// attack-surface pressure over time.
type PressureResult struct {
	Host        string
	AvgOpenPorts float64
	PeakPorts   int
	Pressure    float64 // normalised 0-1 relative to peak
	ScanCount   int
}

// Pressure computes port-pressure scores for every host in the store.
// minScans is the minimum number of scan entries required for a host to
// be included in the results.
func Pressure(store Store, minScans int) []PressureResult {
	hosts := store.Hosts()
	var results []PressureResult

	for _, host := range hosts {
		entries := store.ForHost(host)
		if len(entries) < minScans {
			continue
		}

		var totalPorts int
		peak := 0
		for _, e := range entries {
			n := len(e.Ports)
			totalPorts += n
			if n > peak {
				peak = n
			}
		}

		avg := float64(totalPorts) / float64(len(entries))
		pressure := 0.0
		if peak > 0 {
			pressure = math.Round((avg/float64(peak))*1000) / 1000
		}

		results = append(results, PressureResult{
			Host:        host,
			AvgOpenPorts: math.Round(avg*100) / 100,
			PeakPorts:   peak,
			Pressure:    pressure,
			ScanCount:   len(entries),
		})
	}

	// sort descending by Pressure
	for i := 1; i < len(results); i++ {
		for j := i; j > 0 && results[j].Pressure > results[j-1].Pressure; j-- {
			results[j], results[j-1] = results[j-1], results[j]
		}
	}
	return results
}
