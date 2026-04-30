package history

import (
	"math"
	"sort"
	"time"
)

// CadenceResult holds scan cadence statistics for a single host.
type CadenceResult struct {
	Host        string
	ScanCount   int
	AvgInterval time.Duration
	MinInterval time.Duration
	MaxInterval time.Duration
	Jitter      float64 // stddev of intervals in seconds
	Regular     bool    // true if jitter is below threshold
}

// Cadence analyses the timing regularity of scans for each host.
// minScans is the minimum number of scans required to produce a result.
// jitterThreshold (seconds) determines whether a host is considered regular.
func Cadence(store Store, minScans int, jitterThreshold float64) []CadenceResult {
	hosts := store.Hosts()
	sort.Strings(hosts)

	var results []CadenceResult
	for _, host := range hosts {
		entries := store.ForHost(host)
		if len(entries) < minScans {
			continue
		}

		times := make([]time.Time, len(entries))
		for i, e := range entries {
			times[i] = e.Time
		}
		sort.Slice(times, func(i, j int) bool { return times[i].Before(times[j]) })

		intervals := make([]float64, 0, len(times)-1)
		var minI, maxI time.Duration
		for i := 1; i < len(times); i++ {
			d := times[i].Sub(times[i-1])
			intervals = append(intervals, d.Seconds())
			if i == 1 || d < minI {
				minI = d
			}
			if d > maxI {
				maxI = d
			}
		}

		avg := meanF(intervals)
		jitter := stddevF(intervals, avg)

		results = append(results, CadenceResult{
			Host:        host,
			ScanCount:   len(entries),
			AvgInterval: time.Duration(avg) * time.Second,
			MinInterval: minI,
			MaxInterval: maxI,
			Jitter:      math.Round(jitter*100) / 100,
			Regular:     jitter <= jitterThreshold,
		})
	}
	return results
}
