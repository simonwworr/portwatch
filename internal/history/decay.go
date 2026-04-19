package history

import (
	"math"
	"time"
)

// DecayResult holds the decay score for a host port.
type DecayResult struct {
	Host      string
	Port      int
	LastSeen  time.Time
	AgeDays   float64
	Score     float64 // 1.0 = just seen, approaches 0 over time
}

// Decay computes an exponential decay score for each port seen across all hosts.
// halfLife controls how quickly scores decay (e.g. 7 days means score=0.5 after 7 days).
func Decay(store Store, halfLife time.Duration, now time.Time) []DecayResult {
	type key struct {
		host string
		port int
	}
	lastSeen := map[key]time.Time{}

	for host, entries := range store {
		for _, e := range entries {
			for _, p := range e.Ports {
				k := key{host, p}
				if t, ok := lastSeen[k]; !ok || e.Time.After(t) {
					lastSeen[k] = e.Time
				}
			}
		}
	}

	hl := halfLife.Hours() / 24.0 // in days
	var results []DecayResult
	for k, t := range lastSeen {
		ageDays := now.Sub(t).Hours() / 24.0
		score := math.Pow(0.5, ageDays/hl)
		results = append(results, DecayResult{
			Host:     k.host,
			Port:     k.port,
			LastSeen: t,
			AgeDays:  ageDays,
			Score:    score,
		})
	}
	return results
}
