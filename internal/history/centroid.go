package history

import (
	"math"
	"sort"
)

// CentroidResult holds the computed centroid port profile for a host.
type CentroidResult struct {
	Host     string
	Centroid []int
	Deviation float64
}

// Centroid computes the "average" port set for each host across all scans,
// returning the ports that appear in more than half of scans (majority centroid)
// and a deviation score indicating how much individual scans differ from it.
func Centroid(store Store, minScans int) []CentroidResult {
	hosts := store.Hosts()
	var results []CentroidResult

	for _, host := range hosts {
		entries := store.ForHost(host)
		if len(entries) < minScans {
			continue
		}

		// Count how often each port appears
		freq := map[int]int{}
		for _, e := range entries {
			for _, p := range e.Ports {
				freq[p]++
			}
		}

		// Majority centroid: ports present in >50% of scans
		threshold := float64(len(entries)) * 0.5
		var centroid []int
		for port, count := range freq {
			if float64(count) > threshold {
				centroid = append(centroid, port)
			}
		}
		sort.Ints(centroid)

		centroidSet := toIntSet(centroid)

		// Deviation: average Jaccard distance from centroid across all scans
		var totalDist float64
		for _, e := range entries {
			scanSet := toIntSet(e.Ports)
			totalDist += jaccardDistance(centroidSet, scanSet)
		}
		deviation := totalDist / float64(len(entries))
		deviation = math.Round(deviation*1000) / 1000

		results = append(results, CentroidResult{
			Host:      host,
			Centroid:  centroid,
			Deviation: deviation,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Host < results[j].Host
	})
	return results
}

func jaccardDistance(a, b map[int]struct{}) float64 {
	intersection := 0
	for k := range a {
		if _, ok := b[k]; ok {
			intersection++
		}
	}
	union := len(a) + len(b) - intersection
	if union == 0 {
		return 0
	}
	return 1.0 - float64(intersection)/float64(union)
}
