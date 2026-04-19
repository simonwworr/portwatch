package history

import (
	"math"
	"sort"
	"time"
)

// RhythmResult describes the periodic scan pattern for a host.
type RhythmResult struct {
	Host          string
	ScanCount     int
	AvgIntervalSec float64
	StdDevSec     float64
	Regular       bool // true if stddev < 20% of mean
}

// Rhythm analyses scan timing regularity per host.
func Rhythm(store Store, minScans int) []RhythmResult {
	hosts := store.Hosts()
	sort.Strings(hosts)
	var results []RhythmResult
	for _, host := range hosts {
		entries := store.ForHost(host)
		if len(entries) < minScans {
			continue
		}
		times := make([]time.Time, len(entries))
		for i, e := range entries {
			times[i] = e.Timestamp
		}
		sort.Slice(times, func(i, j int) bool { return times[i].Before(times[j]) })
		var gaps []float64
		for i := 1; i < len(times); i++ {
			gaps = append(gaps, times[i].Sub(times[i-1]).Seconds())
		}
		mean := meanF(gaps)
		std := stddevF(gaps, mean)
		results = append(results, RhythmResult{
			Host:          host,
			ScanCount:     len(entries),
			AvgIntervalSec: mean,
			StdDevSec:     std,
			Regular:       mean > 0 && (std/mean) < 0.20,
		})
	}
	return results
}

func meanF(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	var sum float64
	for _, v := range vals {
		sum += v
	}
	return sum / float64(len(vals))
}

func stddevF(vals []float64, mean float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	var variance float64
	for _, v := range vals {
		d := v - mean
		variance += d * d
	}
	return math.Sqrt(variance / float64(len(vals)))
}
