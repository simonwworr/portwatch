package history

import "math"

// SimilarityResult holds the Jaccard similarity score between two hosts.
type SimilarityResult struct {
	HostA      string
	HostB      string
	Similarity float64
	Shared     []int
	Union      []int
}

// Similarity computes pairwise Jaccard similarity between all hosts
// based on their latest port snapshot. Only pairs with similarity >= minScore
// are returned, sorted descending by score.
func Similarity(store Store, minScore float64) []SimilarityResult {
	hosts := store.Hosts()
	snaps := make(map[string]map[int]struct{}, len(hosts))
	for _, h := range hosts {
		entries := store.ForHost(h)
		if len(entries) == 0 {
			continue
		}
		latest := entries[len(entries)-1]
		set := make(map[int]struct{}, len(latest.Ports))
		for _, p := range latest.Ports {
			set[p] = struct{}{}
		}
		snaps[h] = set
	}

	var results []SimilarityResult
	for i := 0; i < len(hosts); i++ {
		for j := i + 1; j < len(hosts); j++ {
			ha, hb := hosts[i], hosts[j]
			sa, sb := snaps[ha], snaps[hb]
			if len(sa) == 0 && len(sb) == 0 {
				continue
			}
			shared, union := jaccardSets(sa, sb)
			if len(union) == 0 {
				continue
			}
			score := math.Round(float64(len(shared))/float64(len(union))*1000) / 1000
			if score < minScore {
				continue
			}
			results = append(results, SimilarityResult{
				HostA:      ha,
				HostB:      hb,
				Similarity: score,
				Shared:     shared,
				Union:      union,
			})
		}
	}
	sortSimilarity(results)
	return results
}

func jaccardSets(a, b map[int]struct{}) (shared, union []int) {
	seen := make(map[int]struct{})
	for p := range a {
		seen[p] = struct{}{}
		if _, ok := b[p]; ok {
			shared = append(shared, p)
		}
	}
	for p := range b {
		seen[p] = struct{}{}
	}
	for p := range seen {
		union = append(union, p)
	}
	sortInts(shared)
	sortInts(union)
	return
}

func sortSimilarity(results []SimilarityResult) {
	for i := 1; i < len(results); i++ {
		for j := i; j > 0 && results[j].Similarity > results[j-1].Similarity; j-- {
			results[j], results[j-1] = results[j-1], results[j]
		}
	}
}
