package history

import (
	"testing"
	"time"
)

func buildFluxStore() Store {
	now := time.Now()
	st := NewMemoryStore()

	// host-a: ports change every scan → high flux
	st.Append("host-a", Entry{Timestamp: now.Add(-3 * time.Hour), Ports: []int{80, 443}})
	st.Append("host-a", Entry{Timestamp: now.Add(-2 * time.Hour), Ports: []int{22, 8080}})
	st.Append("host-a", Entry{Timestamp: now.Add(-1 * time.Hour), Ports: []int{80, 443}})

	// host-b: ports never change → zero flux
	st.Append("host-b", Entry{Timestamp: now.Add(-3 * time.Hour), Ports: []int{80, 443}})
	st.Append("host-b", Entry{Timestamp: now.Add(-2 * time.Hour), Ports: []int{80, 443}})
	st.Append("host-b", Entry{Timestamp: now.Add(-1 * time.Hour), Ports: []int{80, 443}})

	return st
}

func TestFlux_Basic(t *testing.T) {
	st := buildFluxStore()
	results := Flux(st, 2)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestFlux_HighForChangingHost(t *testing.T) {
	st := buildFluxStore()
	results := Flux(st, 2)
	// host-a should appear first (highest flux)
	if results[0].Host != "host-a" {
		t.Errorf("expected host-a first, got %s", results[0].Host)
	}
	if results[0].Flux <= 0 {
		t.Errorf("expected positive flux for host-a, got %f", results[0].Flux)
	}
}

func TestFlux_ZeroForStableHost(t *testing.T) {
	st := buildFluxStore()
	results := Flux(st, 2)
	var hostB *FluxResult
	for i := range results {
		if results[i].Host == "host-b" {
			hostB = &results[i]
		}
	}
	if hostB == nil {
		t.Fatal("host-b not found in results")
	}
	if hostB.Flux != 0.0 {
		t.Errorf("expected flux=0 for host-b, got %f", hostB.Flux)
	}
}

func TestFlux_MinScansFilters(t *testing.T) {
	st := buildFluxStore()
	results := Flux(st, 10) // require 10 scans — no host qualifies
	if len(results) != 0 {
		t.Errorf("expected 0 results with high minScans, got %d", len(results))
	}
}

func TestFlux_EmptyStore(t *testing.T) {
	st := NewMemoryStore()
	results := Flux(st, 2)
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}
