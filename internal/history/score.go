package history

import "sort"

// HostScore holds a risk score for a host based on port change frequency.
type HostScore struct {
	Host      string  `json:"host"`
	Score     float64 `json:"score"`
	Changes   int     `json:"changes"`
	Openings  int     `json:"openings"`
	Closings  int     `json:"closings"`
}

// Score computes a simple risk score per host from the watch log.
// Score = (openings*2 + closings) / total_changes, so hosts with more
// port openings relative to closings rank higher.
func Score(dir string) ([]HostScore, error) {
	log, err := LoadWatchLog(dir)
	if err != nil {
		return nil, err
	}

	type acc struct {
		openings int
		closings int
	}
	tally := map[string]*acc{}

	for _, ev := range log.Events {
		if _, ok := tally[ev.Host]; !ok {
			tally[ev.Host] = &acc{}
		}
		tally[ev.Host].openings += len(ev.Opened)
		tally[ev.Host].closings += len(ev.Closed)
	}

	results := make([]HostScore, 0, len(tally))
	for host, a := range tally {
		changes := a.openings + a.closings
		score := float64(a.openings*2+a.closings) / float64(max1(changes))
		results = append(results, HostScore{
			Host:     host,
			Score:    score,
			Changes:  changes,
			Openings: a.openings,
			Closings: a.closings,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	return results, nil
}

// TopN returns up to n highest-scored hosts from the provided scores slice.
func TopN(scores []HostScore, n int) []HostScore {
	if n <= 0 || len(scores) == 0 {
		return nil
	}
	if n > len(scores) {
		n = len(scores)
	}
	return scores[:n]
}

func max1(n int) int {
	if n < 1 {
		return 1
	}
	return n
}
