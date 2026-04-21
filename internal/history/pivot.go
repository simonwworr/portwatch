package history

import "sort"

// PivotRow represents a single host's port presence across time buckets.
type PivotRow struct {
	Host    string
	Buckets map[string][]int // bucket label -> open ports
}

// PivotResult holds the full pivot table output.
type PivotResult struct {
	BucketLabels []string
	Rows         []PivotRow
}

// Pivot builds a host × time-bucket matrix of open ports from the store.
// granularity must be "day" or "month".
func Pivot(store Store, granularity string) PivotResult {
	entries := store.All()

	// Collect bucket labels and per-host per-bucket ports.
	bucketSet := map[string]struct{}{}
	// host -> bucket -> port set
	data := map[string]map[string]map[int]struct{}{}

	for _, e := range entries {
		label := bucketKey(e.Time, granularity)
		bucketSet[label] = struct{}{}

		if data[e.Host] == nil {
			data[e.Host] = map[string]map[int]struct{}{}
		}
		if data[e.Host][label] == nil {
			data[e.Host][label] = map[int]struct{}{}
		}
		for _, p := range e.Ports {
			data[e.Host][label][p] = struct{}{}
		}
	}

	// Sort bucket labels.
	labels := make([]string, 0, len(bucketSet))
	for l := range bucketSet {
		labels = append(labels, l)
	}
	sort.Strings(labels)

	// Sort hosts.
	hosts := make([]string, 0, len(data))
	for h := range data {
		hosts = append(hosts, h)
	}
	sort.Strings(hosts)

	rows := make([]PivotRow, 0, len(hosts))
	for _, h := range hosts {
		buckets := map[string][]int{}
		for _, l := range labels {
			portSet := data[h][l]
			ports := make([]int, 0, len(portSet))
			for p := range portSet {
				ports = append(ports, p)
			}
			sort.Ints(ports)
			buckets[l] = ports
		}
		rows = append(rows, PivotRow{Host: h, Buckets: buckets})
	}

	return PivotResult{BucketLabels: labels, Rows: rows}
}
