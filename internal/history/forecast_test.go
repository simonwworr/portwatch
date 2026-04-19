package history

import (
	"testing"
	"time"
)

func buildForecastStore() Store {
	now := time.Now()
	s := Store{}
	// 4 scans for host-a
	for i := 0; i < 4; i++ {
		s.Append("host-a", Entry{
			Timestamp: now.Add(-time.Duration(i) * time.Hour),
			Ports:     []int{80, 443},
		})
	}
	// port 22 only appears once
	s.Append("host-a", Entry{
		Timestamp: now.Add(-5 * time.Hour),
		Ports:     []int{80, 22},
	})
	return s
}

func TestForecast_Basic(t *testing.T) {
	s := buildForecastStore()
	results := Forecast(s, "host-a", 1)
	if len(results) == 0 {
		t.Fatal("expected forecast entries")
	}
	// port 80 should be first (highest probability)
	if results[0].Port != 80 {
		t.Errorf("expected port 80 first, got %d", results[0].Port)
	}
}

func TestForecast_MinCount(t *testing.T) {
	s := buildForecastStore()
	// port 22 seen only once; minCount=2 should exclude it
	results := Forecast(s, "host-a", 2)
	for _, r := range results {
		if r.Port == 22 {
			t.Error("port 22 should be excluded by minCount=2")
		}
	}
}

func TestForecast_Probability(t *testing.T) {
	s := buildForecastStore()
	results := Forecast(s, "host-a", 1)
	for _, r := range results {
		if r.Probability < 0 || r.Probability > 1 {
			t.Errorf("probability out of range: %f", r.Probability)
		}
	}
}

func TestForecast_EmptyStore(t *testing.T) {
	s := Store{}
	results := Forecast(s, "host-a", 1)
	if results != nil {
		t.Error("expected nil for empty store")
	}
}

func TestForecastSince_FiltersOldEntries(t *testing.T) {
	s := buildForecastStore()
	// Only include last 2 hours; port 22 was 5h ago
	since := time.Now().Add(-2 * time.Hour)
	results := ForecastSince(s, "host-a", since, 1)
	for _, r := range results {
		if r.Port == 22 {
			t.Error("port 22 should be excluded by since filter")
		}
	}
}
