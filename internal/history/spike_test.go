package history

import (
	"testing"
	"time"
)

func buildSpikeStore() *Store {
	s := &Store{}
	base := time.Now().Add(-10 * 24 * time.Hour)
	for i := 0; i < 5; i++ {
		s.Append(Entry{
			Host:      "host1",
			Ports:     []int{80, 443},
			ScannedAt: base.Add(time.Duration(i) * 24 * time.Hour),
		})
	}
	// spike entry
	s.Append(Entry{
		Host:      "host1",
		Ports:     []int{80, 443, 8080, 8443, 9000, 9001, 9002, 9003},
		ScannedAt: base.Add(6 * 24 * time.Hour),
	})
	return s
}

func TestDetectSpikes_FindsSpike(t *testing.T) {
	s := buildSpikeStore()
	results := DetectSpikes(s, "host1", 2.0)
	if len(results) == 0 {
		t.Fatal("expected at least one spike")
	}
	if results[0].Host != "host1" {
		t.Errorf("unexpected host: %s", results[0].Host)
	}
	if results[0].PortCount < 4 {
		t.Errorf("expected high port count, got %d", results[0].PortCount)
	}
}

func TestDetectSpikes_NoSpike(t *testing.T) {
	s := &Store{}
	base := time.Now()
	for i := 0; i < 4; i++ {
		s.Append(Entry{
			Host:      "host1",
			Ports:     []int{80, 443},
			ScannedAt: base.Add(time.Duration(i) * 24 * time.Hour),
		})
	}
	results := DetectSpikes(s, "host1", 3.0)
	if len(results) != 0 {
		t.Errorf("expected no spikes, got %d", len(results))
	}
}

func TestDetectSpikes_InsufficientEntries(t *testing.T) {
	s := &Store{}
	s.Append(Entry{Host: "host1", Ports: []int{80}, ScannedAt: time.Now()})
	results := DetectSpikes(s, "host1", 1.5)
	if results != nil {
		t.Error("expected nil for insufficient entries")
	}
}

func TestDetectSpikes_EmptyStore(t *testing.T) {
	s := &Store{}
	results := DetectSpikes(s, "host1", 2.0)
	if results != nil {
		t.Error("expected nil for empty store")
	}
}
