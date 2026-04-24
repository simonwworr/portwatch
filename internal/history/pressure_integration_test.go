package history

import (
	"testing"
	"time"
)

// TestPressure_ThreeHosts verifies ordering and correctness with three hosts
// that have clearly different pressure profiles.
func TestPressure_ThreeHosts(t *testing.T) {
	now := time.Now()
	s := NewMemoryStore()

	// host-x: always 4 ports → pressure 1.0
	for i := 0; i < 3; i++ {
		s.Append(Entry{Host: "host-x", Ports: []int{22, 80, 443, 8080}, Timestamp: now.Add(time.Duration(i) * time.Hour)})
	}

	// host-y: 1 port then 4 ports alternating → avg≈2.5, peak=4, pressure≈0.625
	s.Append(Entry{Host: "host-y", Ports: []int{22}, Timestamp: now})
	s.Append(Entry{Host: "host-y", Ports: []int{22, 80, 443, 8080}, Timestamp: now.Add(time.Hour)})
	s.Append(Entry{Host: "host-y", Ports: []int{22}, Timestamp: now.Add(2 * time.Hour)})
	s.Append(Entry{Host: "host-y", Ports: []int{22, 80, 443, 8080}, Timestamp: now.Add(3 * time.Hour)})

	// host-z: always 1 port → avg=1, peak=1, pressure=1.0 (tied with host-x)
	for i := 0; i < 3; i++ {
		s.Append(Entry{Host: "host-z", Ports: []int{22}, Timestamp: now.Add(time.Duration(i) * time.Hour)})
	}

	res := Pressure(s, 2)
	if len(res) != 3 {
		t.Fatalf("expected 3 results, got %d", len(res))
	}

	// host-y must be last (lowest pressure)
	last := res[len(res)-1]
	if last.Host != "host-y" {
		t.Errorf("expected host-y last, got %s", last.Host)
	}

	// pressure values must be in descending order
	for i := 1; i < len(res); i++ {
		if res[i].Pressure > res[i-1].Pressure {
			t.Errorf("not sorted descending at index %d", i)
		}
	}
}

func TestPressure_ZeroPeakEdgeCase(t *testing.T) {
	now := time.Now()
	s := NewMemoryStore()
	// host with zero ports in every scan
	for i := 0; i < 3; i++ {
		s.Append(Entry{Host: "ghost", Ports: []int{}, Timestamp: now.Add(time.Duration(i) * time.Hour)})
	}
	res := Pressure(s, 2)
	if len(res) != 1 {
		t.Fatalf("expected 1 result, got %d", len(res))
	}
	if res[0].Pressure != 0.0 {
		t.Errorf("expected pressure 0.0 for zero-port host, got %f", res[0].Pressure)
	}
}
