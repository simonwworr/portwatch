package history

import "math"

// OutlierResult describes a host/port combination whose open-port count
// deviates significantly from that host's historical mean.
type OutlierResult struct {
	Host    string
	Port    int
	ZScore  float64
}

// DetectOutliers returns ports whose appearance frequency across scans is
// more than `threshold` standard deviations from the mean frequency for
// that host. A threshold of 2.0 is a reasonable default.
func DetectOutliers(store Store, host string, threshold float64) []OutlierResult {
	entries := store.ForHost(host)
	if len(entries) < 3 {
		return nil
	}

	// Count how many scans each port appeared in.
	freq := map[int]int{}
	for _, e := range entries {
		for _, p := range e.Ports {
			freq[p]++
		}
	}
	if len(freq) == 0 {
		return nil
	}

	// Compute mean and stddev of frequencies.
	var sum float64
	for _, c := range freq {
		sum += float64(c)
	}
	mean := sum / float64(len(freq))

	var variance float64
	for _, c := range freq {
		d := float64(c) - mean
		variance += d * d
	}
	variance /= float64(len(freq))
	stddev := math.Sqrt(variance)
	if stddev == 0 {
		return nil
	}

	var results []OutlierResult
	for port, count := range freq {
		z := math.Abs(float64(count)-mean) / stddev
		if z >= threshold {
			results = append(results, OutlierResult{
				Host:   host,
				Port:   port,
				ZScore: z,
			})
		}
	}
	return results
}
