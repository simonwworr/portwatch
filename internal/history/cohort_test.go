package history

import (
	"testing"
	"time"
)

func buildCohortStore() Store {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	m := NewMemoryStore()
	// Two hosts first seen in Jan 2024 sharing ports 80 and 443
	m.Append("host-a", Entry{Time: base, Ports: []int{80, 443, 8080}})
	m.Append("host-a", Entry{Time: base.Add(24 * time.Hour), Ports: []int{80, 443}})
	m.Append("host-b", Entry{Time: base.Add(2 * time.Hour), Ports: []int{80, 443, 9090}})
	// One host first seen in Feb 2024 — different cohort
	feb := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	m.Append("host-c", Entry{Time: feb, Ports: []int{22, 80}})
	return m
}

func TestCohort_GroupsByMonth(t *testing.T) {
	store := buildCohortStore()
	result := Cohort(store, 1)
	if len(result) != 2 {
		t.Fatalf("expected 2 cohorts, got %d", len(result))
	}
	if result[0].CohortID != "2024-01" {
		t.Errorf("expected first cohort 2024-01, got %s", result[0].CohortID)
	}
	if result[1].CohortID != "2024-02" {
		t.Errorf("expected second cohort 2024-02, got %s", result[1].CohortID)
	}
}

func TestCohort_SharedPorts(t *testing.T) {
	store := buildCohortStore()
	result := Cohort(store, 1)
	janCohort := result[0]
	// shared between host-a and host-b first entries: 80 and 443
	if len(janCohort.Ports) != 2 {
		t.Fatalf("expected 2 shared ports, got %d: %v", len(janCohort.Ports), janCohort.Ports)
	}
	if janCohort.Ports[0] != 80 || janCohort.Ports[1] != 443 {
		t.Errorf("unexpected shared ports: %v", janCohort.Ports)
	}
}

func TestCohort_MinHostsFilters(t *testing.T) {
	store := buildCohortStore()
	// minHosts=2 should exclude the Feb cohort (only 1 host)
	result := Cohort(store, 2)
	if len(result) != 1 {
		t.Fatalf("expected 1 cohort with minHosts=2, got %d", len(result))
	}
	if result[0].CohortID != "2024-01" {
		t.Errorf("unexpected cohort: %s", result[0].CohortID)
	}
}

func TestCohort_EmptyStore(t *testing.T) {
	m := NewMemoryStore()
	result := Cohort(m, 1)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}

func TestCohort_HostsAreSorted(t *testing.T) {
	store := buildCohortStore()
	result := Cohort(store, 1)
	janCohort := result[0]
	if len(janCohort.Hosts) != 2 {
		t.Fatalf("expected 2 hosts in Jan cohort, got %d", len(janCohort.Hosts))
	}
	if janCohort.Hosts[0] != "host-a" || janCohort.Hosts[1] != "host-b" {
		t.Errorf("hosts not sorted: %v", janCohort.Hosts)
	}
}
