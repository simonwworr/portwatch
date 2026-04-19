package history

import (
	"testing"
	"time"
)

func TestPattern_MultipleHosts(t *testing.T) {
	s := NewMemoryStore()
	now := time.Now()

	for i := 0; i < 5; i++ {
		s.Append("alpha", Entry{Timestamp: now.Add(time.Duration(i) * time.Hour), Ports: []int{22, 80}})
	}
	for i := 0; i < 5; i++ {
		ports := []int{443}
		if i%2 == 0 {
			ports = append(ports, 8080)
		}
		s.Append("beta", Entry{Timestamp: now.Add(time.Duration(i) * time.Hour), Ports: ports})
	}

	results := Pattern(s, 3)

	hosts := map[string]bool{}
	for _, r := range results {
		hosts[r.Host] = true
		if r.Frequency <= 0 || r.Frequency > 1 {
			t.Errorf("frequency out of range for %s port %d: %f", r.Host, r.Port, r.Frequency)
		}
	}
	if !hosts["alpha"] {
		t.Error("expected alpha in results")
	}
	if !hosts["beta"] {
		t.Error("expected beta in results")
	}
}

func TestPattern_SortedByFrequencyDesc(t *testing.T) {
	s := NewMemoryStore()
	now := time.Now()
	for i := 0; i < 4; i++ {
		ports := []int{80}
		if i < 1 {
			ports = append(ports, 9999)
		}
		s.Append("host-x", Entry{Timestamp: now.Add(time.Duration(i) * time.Hour), Ports: ports})
	}
	results := Pattern(s, 2)
	var prev float64 = 2.0
	for _, r := range results {
		if r.Host == "host-x" && r.Frequency > prev {
			t.Errorf("results not sorted desc: %f > %f", r.Frequency, prev)
		}
		if r.Host == "host-x" {
			prev = r.Frequency
		}
	}
}
