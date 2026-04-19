package history

import (
	"testing"
	"time"
)

func buildAnomalyStore(t *testing.T) *Store {
	t.Helper()
	dir := t.TempDir()
	s, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	// host-a: port 80 appears 3/4 times (75%), port 9999 appears 1/4 (25%)
	for i := 0; i < 3; i++ {
		s.Append("host-a", Entry{Timestamp: now.Add(time.Duration(i) * time.Minute), Ports: []int{80}})
	}
	s.Append("host-a", Entry{Timestamp: now.Add(4 * time.Minute), Ports: []int{80, 9999}})
	if err := Save(dir, s); err != nil {
		t.Fatal(err)
	}
	return s
}

func TestDetectAnomalies_FindsRarePort(t *testing.T) {
	s := buildAnomalyStore(t)
	results := DetectAnomalies(s, 50.0)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]
	if r.Host != "host-a" {
		t.Errorf("expected host-a, got %s", r.Host)
	}
	if len(r.RarePorts) != 1 || r.RarePorts[0] != 9999 {
		t.Errorf("expected [9999], got %v", r.RarePorts)
	}
}

func TestDetectAnomalies_NoAnomalies(t *testing.T) {
	s := buildAnomalyStore(t)
	// threshold 10% — 25% is not rare
	results := DetectAnomalies(s, 10.0)
	if len(results) != 0 {
		t.Errorf("expected no anomalies, got %d", len(results))
	}
}

func TestDetectAnomalies_EmptyStore(t *testing.T) {
	dir := t.TempDir()
	s, _ := Load(dir)
	results := DetectAnomalies(s, 50.0)
	if len(results) != 0 {
		t.Errorf("expected empty results")
	}
}
