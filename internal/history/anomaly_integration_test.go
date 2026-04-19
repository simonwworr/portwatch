package history_test

import (
	"testing"
	"time"
)

func TestDetectAnomalies_MultipleHosts(t *testing.T) {
	dir := t.TempDir()
	s, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now()
	// host-a: port 443 always present, port 8080 once out of 5
	for i := 0; i < 4; i++ {
		s.Append("host-a", Entry{Timestamp: now.Add(time.Duration(i) * time.Minute), Ports: []int{443}})
	}
	s.Append("host-a", Entry{Timestamp: now.Add(5 * time.Minute), Ports: []int{443, 8080}})

	// host-b: all ports stable
	for i := 0; i < 4; i++ {
		s.Append("host-b", Entry{Timestamp: now.Add(time.Duration(i) * time.Minute), Ports: []int{22, 80}})
	}

	if err := Save(dir, s); err != nil {
		t.Fatal(err)
	}

	results := DetectAnomalies(s, 30.0)

	if len(results) != 1 {
		t.Fatalf("expected 1 anomaly host, got %d", len(results))
	}
	if results[0].Host != "host-a" {
		t.Errorf("expected host-a, got %s", results[0].Host)
	}
	if len(results[0].RarePorts) != 1 || results[0].RarePorts[0] != 8080 {
		t.Errorf("expected [8080], got %v", results[0].RarePorts)
	}
}
