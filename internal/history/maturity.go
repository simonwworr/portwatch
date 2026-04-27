package history

import (
	"math"
	"sort"
	"time"
)

// MaturityResult holds the maturity score for a host.
// A host with a stable, long-lived port profile scores higher.
type MaturityResult struct {
	Host        string
	Score       float64 // 0.0 – 1.0
	AvgPortAge  float64 // average days a port has been continuously open
	StablePorts []int   // ports open in every observed scan
	ScanCount   int
}

// Maturity scores each host by how "mature" its open-port profile is.
// A mature host has ports that have been consistently open for a long time.
// minScans is the minimum number of scan entries required to include a host.
func Maturity(store Store, minScans int, now time.Time) []MaturityResult {
	hosts := store.Hosts()
	var results []MaturityResult

	for _, host := range hosts {
		entries := store.ForHost(host)
		if len(entries) < minScans {
			continue
		}

		// Determine ports open in every scan (stable ports).
		portCount := map[int]int{}
		for _, e := range entries {
			for _, p := range e.Ports {
				portCount[p]++
			}
		}
		var stablePorts []int
		for p, c := range portCount {
			if c == len(entries) {
				stablePorts = append(stablePorts, p)
			}
		}
		sort.Ints(stablePorts)

		// Average port age: days since each stable port was first seen.
		first := entries[0].Time
		spanDays := now.Sub(first).Hours() / 24
		if spanDays < 1 {
			spanDays = 1
		}

		stableRatio := 0.0
		if len(portCount) > 0 {
			stableRatio = float64(len(stablePorts)) / float64(len(portCount))
		}

		// Score: blend of stable-port ratio and log-normalised time span.
		// log2(spanDays+1) / log2(365+1) caps the time component at ~1 year.
		timeScore := math.Log2(spanDays+1) / math.Log2(366)
		if timeScore > 1 {
			timeScore = 1
		}
		score := 0.6*stableRatio + 0.4*timeScore

		results = append(results, MaturityResult{
			Host:        host,
			Score:       math.Round(score*1000) / 1000,
			AvgPortAge:  math.Round(spanDays*10) / 10,
			StablePorts: stablePorts,
			ScanCount:   len(entries),
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	return results
}
