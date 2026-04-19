package history

import (
	"testing"
	"time"
)

func buildFilterStore(t *testing.T) *Store {
	t.Helper()
	dir := t.TempDir()
	s, _ := Load(dir)
	now := time.Now()
	_ = s.Append("host-a", Entry{Host: "host-a", Ports: []int{80, 443}, Time: now.Add(-2 * time.Hour)})
	_ = s.Append("host-a", Entry{Host: "host-a", Ports: []int{22}, Time: now.Add(-30 * time.Minute)})
	_ = s.Append("host-b", Entry{Host: "host-b", Ports: []int{8080}, Time: now.Add(-1 * time.Hour)})
	_ = Save(dir, s)
	return s
}

func TestFilter_ByHost(t *testing.T) {
	s := buildFilterStore(t)
	res := Filter(s, FilterOptions{Host: "host-a"})
	if len(res) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(res))
	}
}

func TestFilter_ByPort(t *testing.T) {
	s := buildFilterStore(t)
	res := Filter(s, FilterOptions{Port: 22})
	if len(res) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(res))
	}
}

func TestFilter_BySince(t *testing.T) {
	s := buildFilterStore(t)
	res := Filter(s, FilterOptions{Since: time.Now().Add(-45 * time.Minute)})
	if len(res) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(res))
	}
}

func TestFilter_ByUntil(t *testing.T) {
	s := buildFilterStore(t)
	res := Filter(s, FilterOptions{Until: time.Now().Add(-90 * time.Minute)})
	if len(res) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(res))
	}
}

func TestFilter_NoMatch(t *testing.T) {
	s := buildFilterStore(t)
	res := Filter(s, FilterOptions{Host: "unknown"})
	if len(res) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(res))
	}
}
