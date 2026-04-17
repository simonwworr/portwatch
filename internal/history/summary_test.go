package history

import (
	"testing"
	"time"
)

func buildStore(entries []Entry) *Store {
	s := &Store{}
	for _, e := range entries {
		s.Append(e)
	}
	return s
}

func TestSummarize_Basic(t *testing.T) {
	now := time.Now()
	s := buildStore([]Entry{
		{Host: "host1", ScannedAt: now.Add(-2 * time.Hour), OpenPorts: []int{80, 443}},
		{Host: "host1", ScannedAt: now.Add(-1 * time.Hour), OpenPorts: []int{80, 443, 8080}},
		{Host: "host1", ScannedAt: now, OpenPorts: []int{80}},
	})

	sum, ok := Summarize(s, "host1")
	if !ok {
		t.Fatal("expected summary")
	}
	if sum.TotalScans != 3 {
		t.Errorf("TotalScans: got %d want 3", sum.TotalScans)
	}
	if sum.MaxOpen != 3 {
		t.Errorf("MaxOpen: got %d want 3", sum.MaxOpen)
	}
	if sum.MinOpen != 1 {
		t.Errorf("MinOpen: got %d want 1", sum.MinOpen)
	}
	if len(sum.UniquePorts) != 3 {
		t.Errorf("UniquePorts len: got %d want 3", len(sum.UniquePorts))
	}
	if !sum.FirstSeen.Equal(now.Add(-2 * time.Hour)) {
		t.Error("unexpected FirstSeen")
	}
}

func TestSummarize_NoEntries(t *testing.T) {
	s := &Store{}
	_, ok := Summarize(s, "ghost")
	if ok {
		t.Error("expected no summary for unknown host")
	}
}

func TestSummarize_UniquePortsSorted(t *testing.T) {
	now := time.Now()
	s := buildStore([]Entry{
		{Host: "h", ScannedAt: now, OpenPorts: []int{9000, 80, 443}},
	})
	sum, _ := Summarize(s, "h")
	prev := -1
	for _, p := range sum.UniquePorts {
		if p <= prev {
			t.Errorf("ports not sorted: %v", sum.UniquePorts)
			break
		}
		prev = p
	}
}
