package history

import (
	"testing"
	"time"
)

func buildExposureStore() Store {
	now := time.Now()
	s := NewMemoryStore()
	// host-a: 4 scans, port 80 open 4/4, port 443 open 2/4
	for i := 0; i < 4; i++ {
		ports := []int{80}
		if i%2 == 0 {
			ports = append(ports, 443)
		}
		_ = s.Append(Entry{Host: "host-a", Ports: ports, Time: now.Add(time.Duration(i) * time.Hour)})
	}
	// host-b: 2 scans, port 22 open 1/2
	for i := 0; i < 2; i++ {
		ports := []int{}
		if i == 0 {
			ports = []int{22}
		}
		_ = s.Append(Entry{Host: "host-b", Ports: ports, Time: now.Add(time.Duration(i) * time.Hour)})
	}
	return s
}

func TestExposure_Basic(t *testing.T) {
	s := buildExposureStore()
	results := Exposure(s, 1)

	// find port 80 on host-a
	var r80 *ExposureResult
	for i := range results {
		if results[i].Host == "host-a" && results[i].Port == 80 {
			r80 = &results[i]
		}
	}
	if r80 == nil {
		t.Fatal("expected result for host-a:80")
	}
	if r80.Exposure != 1.0 {
		t.Errorf("expected exposure 1.0, got %f", r80.Exposure)
	}
}

func TestExposure_PartialPort(t *testing.T) {
	s := buildExposureStore()
	results := Exposure(s, 1)

	var r443 *ExposureResult
	for i := range results {
		if results[i].Host == "host-a" && results[i].Port == 443 {
			r443 = &results[i]
		}
	}
	if r443 == nil {
		t.Fatal("expected result for host-a:443")
	}
	if r443.Exposure != 0.5 {
		t.Errorf("expected exposure 0.5, got %f", r443.Exposure)
	}
}

func TestExposure_MinScansFilters(t *testing.T) {
	s := buildExposureStore()
	// host-b has 2 scans; minScans=3 should exclude it
	results := Exposure(s, 3)
	for _, r := range results {
		if r.Host == "host-b" {
			t.Error("host-b should have been filtered out")
		}
	}
}

func TestExposure_EmptyStore(t *testing.T) {
	s := NewMemoryStore()
	results := Exposure(s, 1)
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}

func TestExposure_SortedDescending(t *testing.T) {
	s := buildExposureStore()
	results := Exposure(s, 1)
	for i := 1; i < len(results); i++ {
		if results[i].Exposure > results[i-1].Exposure {
			t.Errorf("results not sorted descending at index %d", i)
		}
	}
}
