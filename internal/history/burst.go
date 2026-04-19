package history

import "time"

// BurstResult holds the burst detection result for a host.
type BurstResult struct {
	Host      string
	BurstSize int     // max ports opened in a single window
	WindowStart time.Time
	WindowEnd   time.Time
	PeakPorts []int
}

// DetectBursts finds hosts where the number of newly opened ports within a
// sliding window of windowSize scans exceeds threshold.
func DetectBursts(store *Store, host string, windowSize int, threshold int) []BurstResult {
	entries := store.ForHost(host)
	if len(entries) < 2 {
		return nil
	}

	var results []BurstResult

	for i := 1; i+windowSize-1 < len(entries); i++ {
		window := entries[i : i+windowSize]
		basePorts := toPortSet([]Entry{entries[i-1]})
		opened := map[int]struct{}{}
		for _, e := range window {
			for _, p := range e.Ports {
				if _, seen := basePorts[p]; !seen {
					opened[p] = struct{}{}
				}
			}
		}
		if len(opened) >= threshold {
			peaks := make([]int, 0, len(opened))
			for p := range opened {
				peaks = append(peaks, p)
			}
			sortInts(peaks)
			results = append(results, BurstResult{
				Host:        host,
				BurstSize:   len(opened),
				WindowStart: window[0].Time,
				WindowEnd:   window[len(window)-1].Time,
				PeakPorts:   peaks,
			})
		}
	}

	return results
}
