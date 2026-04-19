package history

import "sort"

// ClusterResult groups hosts that share a common set of open ports.
type ClusterResult struct {
	Ports []int
	Hosts []string
}

// Cluster groups hosts by their latest open-port fingerprint.
// Hosts whose port sets are identical are placed in the same cluster.
func Cluster(store Store) []ClusterResult {
	latest := LatestSnapshot(store)

	// key: canonical port string → cluster
	type entry struct {
		ports []int
		hosts []string
	}
	index := map[string]*entry{}
	order := []string{}

	for host, ports := range latest {
		sorted := make([]int, len(ports))
		copy(sorted, ports)
		sort.Ints(sorted)
		key := intsKey(sorted)
		if _, ok := index[key]; !ok {
			index[key] = &entry{ports: sorted}
			order = append(order, key)
		}
		index[key].hosts = append(index[key].hosts, host)
	}

	results := make([]ClusterResult, 0, len(order))
	for _, k := range order {
		e := index[k]
		sort.Strings(e.hosts)
		results = append(results, ClusterResult{Ports: e.ports, Hosts: e.hosts})
	}
	return results
}

// intsKey returns a stable string key for a sorted int slice.
func intsKey(ports []int) string {
	b := make([]byte, 0, len(ports)*4)
	for _, p := range ports {
		b = append(b, byte(p>>24), byte(p>>16), byte(p>>8), byte(p))
	}
	return string(b)
}
