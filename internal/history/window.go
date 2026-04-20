package history

import "time"

// WindowResult holds aggregated stats for a rolling time window.
type WindowResult struct {
	Host      string
	Window    string
	ScanCount int
	AvgPorts  float64
	MinPorts  int
	MaxPorts  int
	Unique    []int
}

// Window computes rolling window statistics for each host over the given
// duration (e.g. 7 days, 24 hours). Only entries within [now-size, now] are
// considered.
func Window(store Store, size time.Duration) []WindowResult {
	all := store.All()
	now := time.Now()
	cutoff := now.Add(-size)

	type agg struct {
		scans  int
		total  int
		min    int
		max    int
		unique map[int]struct{}
	}

	hosts := map[string]*agg{}

	for _, e := range all {
		if e.Time.Before(cutoff) {
			continue
		}
		a, ok := hosts[e.Host]
		if !ok {
			a = &agg{min: -1, unique: map[int]struct{}{}}
			hosts[e.Host] = a
		}
		n := len(e.Ports)
		a.scans++
		a.total += n
		if a.min == -1 || n < a.min {
			a.min = n
		}
		if n > a.max {
			a.max = n
		}
		for _, p := range e.Ports {
			a.unique[p] = struct{}{}
		}
	}

	label := fmtDuration(size)
	var results []WindowResult
	for host, a := range hosts {
		avg := 0.0
		if a.scans > 0 {
			avg = float64(a.total) / float64(a.scans)
		}
		min := a.min
		if min == -1 {
			min = 0
		}
		u := make([]int, 0, len(a.unique))
		for p := range a.unique {
			u = append(u, p)
		}
		sortInts(u)
		results = append(results, WindowResult{
			Host:      host,
			Window:    label,
			ScanCount: a.scans,
			AvgPorts:  avg,
			MinPorts:  min,
			MaxPorts:  a.max,
			Unique:    u,
		})
	}
	return results
}

func fmtDuration(d time.Duration) string {
	if d%(24*time.Hour) == 0 {
		days := int(d.Hours() / 24)
		if days == 1 {
			return "1d"
		}
		return fmt.Sprintf("%dd", days)
	}
	h := int(d.Hours())
	if h > 0 {
		return fmt.Sprintf("%dh", h)
	}
	return d.String()
}
