package history

import (
	"testing"
	"time"
)

func buildPulseStore() Store {
	now := time.Date(2024, 6, 10, 12, 0, 0, 0, time.UTC)
	s := NewMemoryStore()
	// host-a: active 3 out of 4 days
	s.Append(Entry{Host: "host-a", Time: now.AddDate(0, 0, -3), Ports: []int{80, 443}})
	s.Append(Entry{Host: "host-a", Time: now.AddDate(0, 0, -2), Ports: []int{80}})
	s.Append(Entry{Host: "host-a", Time: now.AddDate(0, 0, -1), Ports: []int{80, 8080}})
	// day 0 (today) missing for host-a
	// host-b: active every day
	s.Append(Entry{Host: "host-b", Time: now.AddDate(0, 0, -3), Ports: []int{22}})
	s.Append(Entry{Host: "host-b", Time: now.AddDate(0, 0, -2), Ports: []int{22}})
	s.Append(Entry{Host: "host-b", Time: now.AddDate(0, 0, -1), Ports: []int{22}})
	s.Append(Entry{Host: "host-b", Time: now, Ports: []int{22}})
	return s
}

func TestPulse_Basic(t *testing.T) {
	s := buildPulseStore()
	since := time.Date(2024, 6, 7, 0, 0, 0, 0, time.UTC)
	results := Pulse(s, since)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	// host-b should score 1.0 (4/4 days)
	if results[0].Host != "host-b" {
		t.Errorf("expected host-b first, got %s", results[0].Host)
	}
	if results[0].PulseScore != 1.0 {
		t.Errorf("expected pulse 1.0, got %f", results[0].PulseScore)
	}
}

func TestPulse_ActivaDays(t *testing.T) {
	s := buildPulseStore()
	since := time.Date(2024, 6, 7, 0, 0, 0, 0, time.UTC)
	results := Pulse(s, since)
	var hostA PulseResult
	for _, r := range results {
		if r.Host == "host-a" {
			hostA = r
		}
	}
	if hostA.ActiveDays != 3 {
		t.Errorf("expected 3 active days for host-a, got %d", hostA.ActiveDays)
	}
	if hostA.TotalDays != 4 {
		t.Errorf("expected 4 total days, got %d", hostA.TotalDays)
	}
}

func TestPulse_AvgPorts(t *testing.T) {
	s := buildPulseStore()
	since := time.Date(2024, 6, 7, 0, 0, 0, 0, time.UTC)
	results := Pulse(s, since)
	var hostA PulseResult
	for _, r := range results {
		if r.Host == "host-a" {
			hostA = r
		}
	}
	// ports: 2+1+2 = 5 across 3 scans => 1.67
	if hostA.AvgPorts < 1.6 || hostA.AvgPorts > 1.7 {
		t.Errorf("unexpected avg ports: %f", hostA.AvgPorts)
	}
}

func TestPulse_EmptyStore(t *testing.T) {
	s := NewMemoryStore()
	since := time.Now().AddDate(0, 0, -7)
	results := Pulse(s, since)
	if results != nil {
		t.Errorf("expected nil for empty store, got %v", results)
	}
}

func TestPulse_SinceFilter(t *testing.T) {
	s := buildPulseStore()
	// since = tomorrow: nothing qualifies
	since := time.Now().AddDate(0, 0, 1)
	results := Pulse(s, since)
	if len(results) != 0 {
		t.Errorf("expected 0 results after future since, got %d", len(results))
	}
}
