package history

import (
	"testing"
	"time"
)

func buildConvergenceStore() Store {
	now := time.Now()
	s := NewMemoryStore()
	// host-a: changes twice then stabilises
	s.Append("host-a", Entry{Timestamp: now.Add(-4 * 24 * time.Hour), Ports: []int{80}})
	s.Append("host-a", Entry{Timestamp: now.Add(-3 * 24 * time.Hour), Ports: []int{80, 443}})
	s.Append("host-a", Entry{Timestamp: now.Add(-2 * 24 * time.Hour), Ports: []int{80, 443, 8080}})
	s.Append("host-a", Entry{Timestamp: now.Add(-1 * 24 * time.Hour), Ports: []int{80, 443, 8080}})
	s.Append("host-a", Entry{Timestamp: now, Ports: []int{80, 443, 8080}})
	// host-b: never stabilises
	s.Append("host-b", Entry{Timestamp: now.Add(-3 * 24 * time.Hour), Ports: []int{22}})
	s.Append("host-b", Entry{Timestamp: now.Add(-2 * 24 * time.Hour), Ports: []int{22, 80}})
	s.Append("host-b", Entry{Timestamp: now.Add(-1 * 24 * time.Hour), Ports: []int{22}})
	s.Append("host-b", Entry{Timestamp: now, Ports: []int{22, 443}})
	return s
}

func TestConvergence_Basic(t *testing.T) {
	res := Convergence(buildConvergenceStore())
	if len(res) != 2 {
		t.Fatalf("expected 2 results, got %d", len(res))
	}
}

func TestConvergence_HostAStable(t *testing.T) {
	res := Convergence(buildConvergenceStore())
	var r *ConvergenceResult
	for i := range res {
		if res[i].Host == "host-a" {
			r = &res[i]
		}
	}
	if r == nil {
		t.Fatal("host-a not found")
	}
	if !r.Stable {
		t.Errorf("expected host-a to be stable")
	}
	if r.ChangeRate >= 1.0 {
		t.Errorf("unexpected change rate %.2f", r.ChangeRate)
	}
}

func TestConvergence_HostBUnstable(t *testing.T) {
	res := Convergence(buildConvergenceStore())
	var r *ConvergenceResult
	for i := range res {
		if res[i].Host == "host-b" {
			r = &res[i]
		}
	}
	if r == nil {
		t.Fatal("host-b not found")
	}
	if r.Stable {
		t.Errorf("expected host-b to be unstable")
	}
	if r.ChangeRate == 0 {
		t.Errorf("expected non-zero change rate for host-b")
	}
}

func TestConvergence_EmptyStore(t *testing.T) {
	res := Convergence(NewMemoryStore())
	if len(res) != 0 {
		t.Errorf("expected empty results, got %d", len(res))
	}
}

func TestConvergence_InsufficientEntries(t *testing.T) {
	s := NewMemoryStore()
	s.Append("solo", Entry{Timestamp: time.Now(), Ports: []int{80}})
	res := Convergence(s)
	if len(res) != 0 {
		t.Errorf("expected no result for single-entry host, got %d", len(res))
	}
}
