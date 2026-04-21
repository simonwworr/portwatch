package history

import (
	"math"
	"sort"
	"time"
)

// PulseResult holds the activity pulse score for a host.
type PulseResult struct {
	Host        string
	PulseScore  float64 // 0.0–1.0, higher = more consistently active
	ScanCount   int
	ActiveDays  int
	TotalDays   int
	AvgPorts    float64
}

// Pulse measures how consistently active each host has been over the
// observation window. A host that opens ports on most days scores near 1.0;
// a host that is rarely active scores near 0.0.
func Pulse(store Store, since time.Time) []PulseResult {
	all := store.All()
	type dayKey struct {
		host string
		day  string
	}

	type hostData struct {
		days      map[string]bool
		totalPorts int
		scans     int
	}

	hosts := map[string]*hostData{}
	earliestDay := ""
	latestDay := ""

	for _, e := range all {
		if e.Time.Before(since) {
			continue
		}
		day := e.Time.Format("2006-01-02")
		if earliestDay == "" || day < earliestDay {
			earliestDay = day
		}
		if day > latestDay {
			latestDay = day
		}
		hd, ok := hosts[e.Host]
		if !ok {
			hd = &hostData{days: map[string]bool{}}
			hosts[e.Host] = hd
		}
		hd.days[day] = true
		hd.totalPorts += len(e.Ports)
		hd.scans++
	}

	if earliestDay == "" {
		return nil
	}

	t0, _ := time.Parse("2006-01-02", earliestDay)
	t1, _ := time.Parse("2006-01-02", latestDay)
	totalDays := int(math.Round(t1.Sub(t0).Hours()/24)) + 1

	results := make([]PulseResult, 0, len(hosts))
	for host, hd := range hosts {
		activeDays := len(hd.days)
		score := 0.0
		if totalDays > 0 {
			score = float64(activeDays) / float64(totalDays)
		}
		avgPorts := 0.0
		if hd.scans > 0 {
			avgPorts = float64(hd.totalPorts) / float64(hd.scans)
		}
		results = append(results, PulseResult{
			Host:       host,
			PulseScore: math.Round(score*1000) / 1000,
			ScanCount:  hd.scans,
			ActiveDays: activeDays,
			TotalDays:  totalDays,
			AvgPorts:   math.Round(avgPorts*100) / 100,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].PulseScore != results[j].PulseScore {
			return results[i].PulseScore > results[j].PulseScore
		}
		return results[i].Host < results[j].Host
	})
	return results
}
