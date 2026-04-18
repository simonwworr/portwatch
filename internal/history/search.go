package history

import "time"

// SearchQuery defines filters for searching history entries.
type SearchQuery struct {
	Host  string
	Port  int
	Since time.Time
	Until time.Time
}

// SearchResult holds a matching entry with its host.
type SearchResult struct {
	Host  string
	Entry Entry
}

// Search returns entries matching the given query across all hosts.
func Search(store *Store, q SearchQuery) []SearchResult {
	var results []SearchResult

	hosts := store.Hosts()
	for _, host := range hosts {
		if q.Host != "" && host != q.Host {
			continue
		}
		entries := store.ForHost(host)
		for _, e := range entries {
			if !q.Since.IsZero() && e.Timestamp.Before(q.Since) {
				continue
			}
			if !q.Until.IsZero() && e.Timestamp.After(q.Until) {
				continue
			}
			if q.Port != 0 && !containsPort(e.Ports, q.Port) {
				continue
			}
			results = append(results, SearchResult{Host: host, Entry: e})
		}
	}
	return results
}

func containsPort(ports []int, port int) bool {
	for _, p := range ports {
		if p == port {
			return true
		}
	}
	return false
}
