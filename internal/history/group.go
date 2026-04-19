package history

import (
	"fmt"
	"sort"
	"time"
)

// GroupEntry holds aggregated scan results for a time bucket.
type GroupEntry struct {
	Bucket      string `json:"bucket"`
	Host        string `json:"host"`
	ScanCount   int    `json:"scan_count"`
	UniquePorts []int  `json:"unique_ports"`
	AvgPorts    float64 `json:"avg_ports"`
}

// GroupBy aggregates history entries for a host by a time bucket: "day", "week", or "month".
func GroupBy(store *Store, host, bucket string) ([]GroupEntry, error) {
	entries := store.ForHost(host)
	if len(entries) == 0 {
		return nil, nil
	}

	type bucketData struct {
		ports      map[int]struct{}
		totalPorts int
		scans      int
	}

	buckets := map[string]*bucketData{}
	keys := []string{}

	for _, e := range entries {
		key := bucketKey(e.Time, bucket)
		if _, ok := buckets[key]; !ok {
			buckets[key] = &bucketData{ports: map[int]struct{}{}}
			keys = append(keys, key)
		}
		d := buckets[key]
		d.scans++
		d.totalPorts += len(e.Ports)
		for _, p := range e.Ports {
			d.ports[p] = struct{}{}
		}
	}

	sort.Strings(keys)
	result := make([]GroupEntry, 0, len(keys))
	for _, k := range keys {
		d := buckets[k]
		uniq := make([]int, 0, len(d.ports))
		for p := range d.ports {
			uniq = append(uniq, p)
		}
		sort.Ints(uniq)
		avg := 0.0
		if d.scans > 0 {
			avg = float64(d.totalPorts) / float64(d.scans)
		}
		result = append(result, GroupEntry{
			Bucket:      k,
			Host:        host,
			ScanCount:   d.scans,
			UniquePorts: uniq,
			AvgPorts:    avg,
		})
	}
	return result, nil
}

func bucketKey(t time.Time, bucket string) string {
	switch bucket {
	case "week":
		year, week := t.ISOWeek()
		return fmt.Sprintf("%d-W%02d", year, week)
	case "month":
		return t.Format("2006-01")
	default:
		return t.Format("2006-01-02")
	}
}
