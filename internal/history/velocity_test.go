package history

import (
	"testing"
	"time"
)

func buildVelocityStore(t *testing.T) *Store {
	t.Helper()
	dir := t.TempDir()
	s, err := Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	now := time.Now()
	s.Append("hostA", Entry{Host: "hostA", Time: now.Add(-6 * 24 * time.Hour), Ports: []int{80}})
	s.Append("hostA", Entry{Host: "hostA", Time: now.Add(-3 * 24 * time.Hour), Ports: []int{80, 443}})
	s.Append("hostA", Entry{Host: "hostA", Time: now.Add(-1 * 24 * time.Hour), Ports: []int{80, 443, 8080}})
	s.Append("hostB", Entry{Host: "hostB", Time: now.Add(-5 * 24 * time.Hour), Ports: []int{22, 80}})
	s.Append("hostB", Entry{Host: "hostB", Time: now.Add(-2 * 24 * time.Hour), Ports: []int{22}})
	return s
}

func TestVelocity_Basic(t *testing.T) {
	s := buildVelocityStore(t)
	results := Velocity(s, 30*24*time.Hour)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestVelocity_OpenRate(t *testing.T) {
	s := buildVelocityStore(t)
	results := Velocity(s, 30*24*time.Hour)
	var hostA VelocityEntry
	for _, r := range results {
		if r.Host == "hostA" {
			hostA = r
		}
	}
	if hostA.OpenRate <= 0 {
		t.Errorf("expected positive open rate for hostA, got %f", hostA.OpenRate)
	}
	if hostA.CloseRate != 0 {
		t.Errorf("expected zero close rate for hostA, got %f", hostA.CloseRate)
	}
}

func TestVelocity_CloseRate(t *testing.T) {
	s := buildVelocityStore(t)
	results := Velocity(s, 30*24*time.Hour)
	var hostB VelocityEntry
	for _, r := range results {
		if r.Host == "hostB" {
			hostB = r
		}
	}
	if hostB.CloseRate <= 0 {
		t.Errorf("expected positive close rate for hostB, got %f", hostB.CloseRate)
	}
}

func TestVelocity_InsufficientEntries(t *testing.T) {
	dir := t.TempDir()
	s, _ := Load(dir)
	s.Append("lonely", Entry{Host: "lonely", Time: time.Now(), Ports: []int{80}})
	results := Velocity(s, 7*24*time.Hour)
	if len(results) != 0 {
		t.Errorf("expected 0 results for single-entry host, got %d", len(results))
	}
}

func TestVelocity_EmptyStore(t *testing.T) {
	dir := t.TempDir()
	s, _ := Load(dir)
	results := Velocity(s, 7*24*time.Hour)
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}
