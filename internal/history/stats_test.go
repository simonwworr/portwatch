package history

import (
	"testing"
	"time"
)

func buildStatsStore() *Store {
	s := &Store{entries: make(map[string][]Entry)}
	now := time.Now()
	s.entries["host-a"] = []Entry{
		{Host: "host-a", Ports: []int{80, 443, 22}, ScannedAt: now.Add(-2 * time.Hour)},
		{Host: "host-a", Ports: []int{80, 443}, ScannedAt: now.Add(-1 * time.Hour)},
		{Host: "host-a", Ports: []int{80}, ScannedAt: now},
	}
	s.entries["host-b"] = []Entry{
		{Host: "host-b", Ports: []int{8080}, ScannedAt: now},
	}
	return s
}

func TestStats_ScanCount(t *testing.T) {
	s := buildStatsStore()
	result := Stats(s, 0)
	if len(result) != 2 {
		t.Fatalf("expected 2 hosts, got %d", len(result))
	}
	if result[0].Host != "host-a" {
		t.Fatalf("expected host-a first, got %s", result[0].Host)
	}
	if result[0].ScanCount != 3 {
		t.Errorf("expected 3 scans for host-a, got %d", result[0].ScanCount)
	}
	if result[1].ScanCount != 1 {
		t.Errorf("expected 1 scan for host-b, got %d", result[1].ScanCount)
	}
}

func TestStats_PortFrequency(t *testing.T) {
	s := buildStatsStore()
	result := Stats(s, 0)
	top := result[0].TopPorts // host-a
	if top[0].Port != 80 || top[0].SeenCount != 3 {
		t.Errorf("expected port 80 seen 3 times, got port %d seen %d", top[0].Port, top[0].SeenCount)
	}
	if top[1].SeenCount != 2 {
		t.Errorf("expected second port seen 2 times, got %d", top[1].SeenCount)
	}
}

func TestStats_TopN(t *testing.T) {
	s := buildStatsStore()
	result := Stats(s, 1)
	if len(result[0].TopPorts) != 1 {
		t.Errorf("expected 1 top port, got %d", len(result[0].TopPorts))
	}
}

func TestStats_EmptyStore(t *testing.T) {
	s := &Store{entries: make(map[string][]Entry)}
	result := Stats(s, 5)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d", len(result))
	}
}
