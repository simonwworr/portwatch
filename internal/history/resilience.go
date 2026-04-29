package history

import "math"

// ResilienceResult holds the resilience score for a single host.
// Resilience measures how quickly a host recovers its port set after a
// disruption (port loss followed by port regain).
type ResilienceResult struct {
	Host            string  `json:"host"`
	Scans           int     `json:"scans"`
	Disruptions     int     `json:"disruptions"`
	AvgRecoveryRate float64 `json:"avg_recovery_rate"` // 0.0–1.0
	Score           float64 `json:"score"`
}

// Resilience analyses the scan history for each host and returns a
// ResilienceResult per host. A disruption is defined as a scan where the
// number of open ports drops by more than dropThreshold fraction compared to
// the previous scan. The recovery rate is the fraction of lost ports that
// return in the immediately following scan.
func Resilience(store Store, minScans int, dropThreshold float64) []ResilienceResult {
	hosts := store.Hosts()
	var results []ResilienceResult

	for _, host := range hosts {
		entries := store.ForHost(host)
		if len(entries) < minScans {
			continue
		}

		disruptions := 0
		totalRecovery := 0.0

		for i := 1; i < len(entries); i++ {
			prev := toIntSet(entries[i-1].Ports)
			curr := toIntSet(entries[i].Ports)

			prevCount := len(prev)
			if prevCount == 0 {
				continue
			}

			// Count ports lost in this scan.
			lost := 0
			for p := range prev {
				if _, ok := curr[p]; !ok {
					lost++
				}
			}

			dropFrac := float64(lost) / float64(prevCount)
			if dropFrac < dropThreshold {
				continue
			}

			disruptions++

			// Measure recovery in the next scan if available.
			if i+1 >= len(entries) {
				// No following scan; assume no recovery.
				continue
			}

			next := toIntSet(entries[i+1].Ports)
			recovered := 0
			for p := range prev {
				if _, ok := curr[p]; !ok {
					if _, ok2 := next[p]; ok2 {
						recovered++
					}
				}
			}

			if lost > 0 {
				totalRecovery += float64(recovered) / float64(lost)
			}
		}

		avgRecovery := 0.0
		if disruptions > 0 {
			avgRecovery = totalRecovery / float64(disruptions)
		}

		// Score combines low disruption rate with high recovery rate.
		disruptionRate := float64(disruptions) / float64(len(entries)-1)
		score := math.Round((avgRecovery*(1.0-disruptionRate))*100) / 100

		results = append(results, ResilienceResult{
			Host:            host,
			Scans:           len(entries),
			Disruptions:     disruptions,
			AvgRecoveryRate: math.Round(avgRecovery*100) / 100,
			Score:           score,
		})
	}

	return results
}
