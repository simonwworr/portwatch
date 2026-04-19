package history

import "time"

// SpikeResult represents a detected port count spike for a host.
type SpikeResult struct {
	Host      string
	At        time.Time
	PortCount int
	Avg       float64
	Delta     float64
}

// DetectSpikes finds scan entries where the open port count exceeds avg by threshold factor.
func DetectSpikes(store *Store, host string, threshold float64) []SpikeResult {
	entries := store.ForHost(host)
	if len(entries) < 2 {
		return nil
	}

	var sum float64
	for _, e := range entries {
		sum += float64(len(e.Ports))
	}
	avg := sum / float64(len(entries))

	var spikes []SpikeResult
	for _, e := range entries {
		count := float64(len(e.Ports))
		if avg > 0 && count/avg >= threshold {
			spikes = append(spikes, SpikeResult{
				Host:      host,
				At:        e.ScannedAt,
				PortCount: len(e.Ports),
				Avg:       avg,
				Delta:     count - avg,
			})
		}
	}
	return spikes
}
