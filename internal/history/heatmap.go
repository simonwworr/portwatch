package history

import "sort"

// HeatmapEntry represents port activity count for a host over a time bucket.
type HeatmapEntry struct {
	Host   string
	Bucket string // e.g. "2024-01-15"
	Count  int    // number of distinct ports seen
}

// Heatmap builds a day-bucketed activity heatmap from the store.
// Each entry records how many distinct ports were open for a host on a given day.
func Heatmap(store Store) []HeatmapEntry {
	type key struct{ host, bucket string }
	portSets := map[key]map[int]struct{}{}

	for host, entries := range store {
		for _, e := range entries {
			bucket := e.Time.UTC().Format("2006-01-02")
			k := key{host, bucket}
			if portSets[k] == nil {
				portSets[k] = map[int]struct{}{}
			}
			for _, p := range e.Ports {
				portSets[k][p] = struct{}{}
			}
		}
	}

	result := make([]HeatmapEntry, 0, len(portSets))
	for k, ports := range portSets {
		result = append(result, HeatmapEntry{
			Host:   k.host,
			Bucket: k.bucket,
			Count:  len(ports),
		})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Host != result[j].Host {
			return result[i].Host < result[j].Host
		}
		return result[i].Bucket < result[j].Bucket
	})
	return result
}
