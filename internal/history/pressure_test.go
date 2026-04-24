package history

import (
	"testing"
	"time"
)

func buildPressureStore() Store {
	now := time.Now()
	s := NewMemoryStore()
	// host-a: always 2 ports open → avg=2, peak=2, pressure=1.0
	for i := 0; i < 4; i++ {
		s.Append(Entry{Host: "host-a", Ports: []int{80, 443}, Timestamp: now.Add(time.Duration(i) * time.Hour)})
	}
	// host-b: oscillates 1 and 3 ports → avg=2, peak=3, pressure≈0.667
	s.Append(Entry{Host: "host-b", Ports: []int{22}, Timestamp: now})
	s.Append(Entry{Host: "host-b", Ports: []int{22, 80, 443}, Timestamp: now.Add(time.Hour)})
	s.Append(Entry{Host: "host-b", Ports: []int{22}, Timestamp: now.Add(2 * time.Hour)})
	s.Append(Entry{Host: "host-b", Ports: []int{22, 80, 443}, Timestamp: now.Add(3 * time.Hour)})
	// host-c: only 1 scan → filtered by minScans=2
	s.Append(Entry{Host: "host-c", Ports: []int{8080}, Timestamp: now})
	return s
}

func TestPressure_Basic(t *testing.T) {
	s := buildPressureStore()
	res := Pressure(s, 2)
	if len(res) != 2 {
		t.Fatalf("expected 2 results, got %d", len(res))
	}
}

func TestPressure_HostA_ScoreOne(t *testing.T) {
	s := buildPressureStore()
	res := Pressure(s, 2)
	// host-a should score 1.0 (avg == peak)
	var found bool
	for _, r := range res {
		if r.Host == "host-a" {
			found = true
			if r.Pressure != 1.0 {
				t.Errorf("host-a pressure: want 1.0, got %f", r.Pressure)
			}
			if r.PeakPorts != 2 {
				t.Errorf("host-a peak: want 2, got %d", r.PeakPorts)
			}
		}
	}
	if !found {
		t.Fatal("host-a not found in results")
	}
}

func TestPressure_HostB_LowerPressure(t *testing.T) {
	s := buildPressureStore()
	res := Pressure(s, 2)
	for _, r := range res {
		if r.Host == "host-b" {
			if r.Pressure >= 1.0 {
				t.Errorf("host-b pressure should be < 1.0, got %f", r.Pressure)
			}
			if r.PeakPorts != 3 {
				t.Errorf("host-b peak: want 3, got %d", r.PeakPorts)
			}
			return
		}
	}
	t.Fatal("host-b not found in results")
}

func TestPressure_MinScansFilters(t *testing.T) {
	s := buildPressureStore()
	res := Pressure(s, 2)
	for _, r := range res {
		if r.Host == "host-c" {
			t.Error("host-c should have been filtered out")
		}
	}
}

func TestPressure_EmptyStore(t *testing.T) {
	s := NewMemoryStore()
	res := Pressure(s, 1)
	if len(res) != 0 {
		t.Errorf("expected empty result, got %d", len(res))
	}
}

func TestPressure_SortedDescending(t *testing.T) {
	s := buildPressureStore()
	res := Pressure(s, 2)
	for i := 1; i < len(res); i++ {
		if res[i].Pressure > res[i-1].Pressure {
			t.Errorf("results not sorted descending at index %d", i)
		}
	}
}
