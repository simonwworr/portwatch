package history

import (
	"fmt"
	"sort"
	"strings"
)

// Fingerprint represents a port-set signature for a host at a point in time.
type Fingerprint struct {
	Host      string
	Signature string
	Ports     []int
}

// Fingerprints computes a deterministic signature for each host based on its
// latest open ports. Hosts with identical open-port sets share a signature.
func Fingerprints(store Store) []Fingerprint {
	snaps := AllSnapshots(store)

	results := make([]Fingerprint, 0, len(snaps))
	for host, entry := range snaps {
		ports := make([]int, len(entry.Ports))
		copy(ports, entry.Ports)
		sort.Ints(ports)

		parts := make([]string, len(ports))
		for i, p := range ports {
			parts[i] = fmt.Sprintf("%d", p)
		}
		sig := strings.Join(parts, ",")
		if sig == "" {
			sig = "<empty>"
		}

		results = append(results, Fingerprint{
			Host:      host,
			Signature: sig,
			Ports:     ports,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Host < results[j].Host
	})
	return results
}

// GroupByFingerprint groups hosts that share the same port-set signature.
func GroupByFingerprint(store Store) map[string][]string {
	fingerprints := Fingerprints(store)
	groups := make(map[string][]string)
	for _, f := range fingerprints {
		groups[f.Signature] = append(groups[f.Signature], f.Host)
	}
	return groups
}
