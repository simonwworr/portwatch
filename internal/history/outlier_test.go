package history

import (
	"testing"
	"time"
)

func buildOutlierStore() Store {
	s := NewMemoryStore()
	now := time.Now()
	// port 80 appears in every scan (high frequency)
	// port 9999 appears only once (low frequency)
	// ports 443, 8080 appear in most scans
	scans := [][]int{
		{80, 443, 8080},
		{80, 443, 8080},
		{80, 443, 8080},
		{80, 443, 8080},
		{80, 443, 8080, 9999},
	}
	for i, ports := range scans {
		s.Append("host1", Entry{
			Timestamp: now.Add(time.Duration(i) * time.Hour),
			Ports:     ports,
		})
	}
	return s
}

func TestDetectOutliers_FindsRarePort(t *testing.T) {
	s := buildOutlierStore()
	results := DetectOutliers(s, "host1", 1.5)
	if len(results) == 0 {
		t.Fatal("expected at least one outlier")
	}
	found := false
	for _, r := range results {
		if r.Port == 9999 {
			found = true
			if r.ZScore <= 0 {
				t.Errorf("expected positive z-score, got %f", r.ZScore)
			}
		}
	}
	if !found {
		t.Error("expected port 9999 to be flagged as outlier")
	}
}

func TestDetectOutliers_NoAnomalies(t *testing.T) {
	s := NewMemoryStore()
	now := time.Now()
	for i := 0; i < 5; i++ {
		s.Append("host1", Entry{Timestamp: now.Add(time.Duration(i) * time.Hour), Ports: []int{80, 443}})
	}
	results := DetectOutliers(s, "host1", 2.0)
	if len(results) != 0 {
		t.Errorf("expected no outliers, got %d", len(results))
	}
}

func TestDetectOutliers_InsufficientEntries(t *testing.T) {
	s := NewMemoryStore()
	s.Append("host1", Entry{Timestamp: time.Now(), Ports: []int{80}})
	results := DetectOutliers(s, "host1", 2.0)
	if results != nil {
		t.Error("expected nil for insufficient entries")
	}
}

func TestDetectOutliers_MissingHost(t *testing.T) {
	s := buildOutlierStore()
	results := DetectOutliers(s, "ghost", 1.5)
	if len(results) != 0 {
		t.Errorf("expected no results for unknown host, got %d", len(results))
	}
}
