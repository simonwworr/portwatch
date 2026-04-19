package history

import (
	"math"
	"testing"
	"time"
)

func buildDecayStore() Store {
	now := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	return Store{
		"host1": {
			{Time: now.Add(-14 * 24 * time.Hour), Ports: []int{80}},
			{Time: now.Add(-1 * 24 * time.Hour), Ports: []int{80, 443}},
		},
		"host2": {
			{Time: now.Add(-7 * 24 * time.Hour), Ports: []int{22}},
		},
	}
}

func TestDecay_ScoreRecentPort(t *testing.T) {
	now := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	store := buildDecayStore()
	results := Decay(store, 7*24*time.Hour, now)

	for _, r := range results {
		if r.Host == "host1" && r.Port == 443 {
			if r.AgeDays < 0.9 || r.AgeDays > 1.1 {
				t.Errorf("expected ~1 day age, got %f", r.AgeDays)
			}
			if r.Score < 0.9 {
				t.Errorf("expected high score for recent port, got %f", r.Score)
			}
		}
	}
}

func TestDecay_HalfLifeAt7Days(t *testing.T) {
	now := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	store := buildDecayStore()
	results := Decay(store, 7*24*time.Hour, now)

	for _, r := range results {
		if r.Host == "host2" && r.Port == 22 {
			if math.Abs(r.Score-0.5) > 0.01 {
				t.Errorf("expected score ~0.5 at half-life, got %f", r.Score)
			}
		}
	}
}

func TestDecay_UsesLatestSeen(t *testing.T) {
	now := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	store := buildDecayStore()
	results := Decay(store, 7*24*time.Hour, now)

	for _, r := range results {
		if r.Host == "host1" && r.Port == 80 {
			// latest seen is 1 day ago, not 14
			if r.AgeDays > 2 {
				t.Errorf("expected latest seen used, age=%f", r.AgeDays)
			}
		}
	}
}

func TestDecay_EmptyStore(t *testing.T) {
	now := time.Now()
	results := Decay(Store{}, 7*24*time.Hour, now)
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}
