package history

import (
	"testing"
	"time"
)

func buildPatternStore() Store {
	s := NewMemoryStore()
	now := time.Now()
	// host-a: 4 scans, port 80 always open, port 443 twice
	for i := 0; i < 4; i++ {
		ports := []int{80}
		if i < 2 {
			ports = append(ports, 443)
		}
		s.Append("host-a", Entry{Timestamp: now.Add(time.Duration(i) * time.Hour), Ports: ports})
	}
	// host-b: 1 scan only (below minScans=2)
	s.Append("host-b", Entry{Timestamp: now, Ports: []int{22}})
	return s
}

func TestPattern_Frequency(t *testing.T) {
	s := buildPatternStore()
	results := Pattern(s, 2)
	if len(results) == 0 {
		t.Fatal("expected results for host-a")
	}
	for _, r := range results {
		if r.Host == "host-b" {
			t.Error("host-b should be excluded (minScans=2)")
		}
	}
}

func TestPattern_Port80AlwaysPresent(t *testing.T) {
	s := buildPatternStore()
	results := Pattern(s, 2)
	var p80 *PortPattern
	for i := range results {
		if results[i].Host == "host-a" && results[i].Port == 80 {
			p80 = &results[i]
			break
		}
	}
	if p80 == nil {
		t.Fatal("port 80 not found")
	}
	if p80.Frequency != 1.0 {
		t.Errorf("expected frequency 1.0, got %f", p80.Frequency)
	}
}

func TestPattern_Port443HalfFrequency(t *testing.T) {
	s := buildPatternStore()
	results := Pattern(s, 2)
	for _, r := range results {
		if r.Host == "host-a" && r.Port == 443 {
			if r.Frequency != 0.5 {
				t.Errorf("expected 0.5, got %f", r.Frequency)
			}
			return
		}
	}
	t.Error("port 443 not found")
}

func TestPattern_EmptyStore(t *testing.T) {
	s := NewMemoryStore()
	results := Pattern(s, 1)
	if len(results) != 0 {
		t.Errorf("expected empty, got %d", len(results))
	}
}
