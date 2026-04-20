package history

import (
	"testing"
	"time"
)

func buildWindowStore() Store {
	now := time.Now()
	s := NewMemoryStore()
	_ = s.Append(Entry{Host: "a", Time: now.Add(-1 * time.Hour), Ports: []int{80, 443}})
	_ = s.Append(Entry{Host: "a", Time: now.Add(-2 * time.Hour), Ports: []int{80}})
	_ = s.Append(Entry{Host: "a", Time: now.Add(-200 * time.Hour), Ports: []int{22}})
	_ = s.Append(Entry{Host: "b", Time: now.Add(-3 * time.Hour), Ports: []int{8080, 9090, 443}})
	return s
}

func TestWindow_ExcludesOldEntries(t *testing.T) {
	s := buildWindowStore()
	results := Window(s, 24*time.Hour)

	var a, b *WindowResult
	for i := range results {
		switch results[i].Host {
		case "a":
			a = &results[i]
		case "b":
			b = &results[i]
		}
	}

	if a == nil {
		t.Fatal("expected result for host a")
	}
	if a.ScanCount != 2 {
		t.Errorf("expected 2 scans for a, got %d", a.ScanCount)
	}
	if b == nil {
		t.Fatal("expected result for host b")
	}
	if b.ScanCount != 1 {
		t.Errorf("expected 1 scan for b, got %d", b.ScanCount)
	}
}

func TestWindow_AvgPorts(t *testing.T) {
	s := buildWindowStore()
	results := Window(s, 24*time.Hour)

	for _, r := range results {
		if r.Host == "a" {
			// scans: [80,443]=2 and [80]=1 → avg 1.5
			if r.AvgPorts != 1.5 {
				t.Errorf("expected avg 1.5, got %f", r.AvgPorts)
			}
		}
	}
}

func TestWindow_UniquePortsSorted(t *testing.T) {
	s := buildWindowStore()
	results := Window(s, 24*time.Hour)

	for _, r := range results {
		if r.Host == "a" {
			if len(r.Unique) != 2 {
				t.Fatalf("expected 2 unique ports, got %d", len(r.Unique))
			}
			if r.Unique[0] != 80 || r.Unique[1] != 443 {
				t.Errorf("unexpected unique ports: %v", r.Unique)
			}
		}
	}
}

func TestWindow_EmptyStore(t *testing.T) {
	s := NewMemoryStore()
	results := Window(s, 24*time.Hour)
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}

func TestWindow_Label(t *testing.T) {
	s := buildWindowStore()
	results := Window(s, 7*24*time.Hour)
	for _, r := range results {
		if r.Window != "7d" {
			t.Errorf("expected window label '7d', got %q", r.Window)
		}
	}
}
