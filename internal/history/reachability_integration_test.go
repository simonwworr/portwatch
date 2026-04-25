package history

import (
	"testing"
	"time"
)

func TestReachability_SortedByUptimeDesc(t *testing.T) {
	now := time.Now()
	s := NewMemoryStore()

	// hostX: 2/4 scans active (uptime 0.5)
	for i := 0; i < 4; i++ {
		ports := []int{}
		if i%2 == 0 {
			ports = []int{80}
		}
		s.Append("hostX", Entry{Timestamp: now.Add(time.Duration(-i) * time.Hour), Ports: ports})
	}

	// hostY: 4/4 scans active (uptime 1.0)
	for i := 0; i < 4; i++ {
		s.Append("hostY", Entry{Timestamp: now.Add(time.Duration(-i) * time.Hour), Ports: []int{443}})
	}

	// hostZ: 1/4 scans active (uptime 0.25)
	for i := 0; i < 4; i++ {
		ports := []int{}
		if i == 0 {
			ports = []int{22}
		}
		s.Append("hostZ", Entry{Timestamp: now.Add(time.Duration(-i) * time.Hour), Ports: ports})
	}

	res := Reachability(s, 1)
	if len(res) != 3 {
		t.Fatalf("expected 3 results, got %d", len(res))
	}

	if res[0].Host != "hostY" {
		t.Errorf("expected hostY first (uptime=1.0), got %s", res[0].Host)
	}
	if res[1].Host != "hostX" {
		t.Errorf("expected hostX second (uptime=0.5), got %s", res[1].Host)
	}
	if res[2].Host != "hostZ" {
		t.Errorf("expected hostZ third (uptime=0.25), got %s", res[2].Host)
	}
}

func TestReachability_AvgPortsCorrect(t *testing.T) {
	now := time.Now()
	s := NewMemoryStore()
	s.Append("h1", Entry{Timestamp: now.Add(-2 * time.Hour), Ports: []int{80, 443, 22}})
	s.Append("h1", Entry{Timestamp: now.Add(-1 * time.Hour), Ports: []int{80}})
	s.Append("h1", Entry{Timestamp: now, Ports: []int{80, 443}})

	res := Reachability(s, 1)
	if len(res) != 1 {
		t.Fatalf("expected 1 result, got %d", len(res))
	}
	// avg = (3+1+2)/3 = 2.0
	expected := 2.0
	if res[0].AvgPorts < expected-0.001 || res[0].AvgPorts > expected+0.001 {
		t.Errorf("expected AvgPorts=%.2f, got %.2f", expected, res[0].AvgPorts)
	}
}
