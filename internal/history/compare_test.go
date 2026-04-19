package history

import (
	"testing"
	"time"
)

func buildCompareStore(t *testing.T) *Store {
	t.Helper()
	store := &Store{}
	now := time.Now()
	store.Append(Entry{Host: "host1", Ports: []int{80, 443}, ScannedAt: now.Add(-2 * time.Minute)})
	store.Append(Entry{Host: "host1", Ports: []int{443, 8080}, ScannedAt: now})
	store.Append(Entry{Host: "host2", Ports: []int{22}, ScannedAt: now})
	return store
}

func TestCompare_OpenedAndClosed(t *testing.T) {
	store := buildCompareStore(t)
	deltas := Compare(store)

	if len(deltas) != 1 {
		t.Fatalf("expected 1 delta (host2 skipped), got %d", len(deltas))
	}
	d := deltas[0]
	if d.Host != "host1" {
		t.Errorf("expected host1, got %s", d.Host)
	}
	if len(d.Opened) != 1 || d.Opened[0] != 8080 {
		t.Errorf("expected opened=[8080], got %v", d.Opened)
	}
	if len(d.Closed) != 1 || d.Closed[0] != 80 {
		t.Errorf("expected closed=[80], got %v", d.Closed)
	}
	if len(d.Stable) != 1 || d.Stable[0] != 443 {
		t.Errorf("expected stable=[443], got %v", d.Stable)
	}
}

func TestCompare_NoChanges(t *testing.T) {
	store := &Store{}
	now := time.Now()
	store.Append(Entry{Host: "h", Ports: []int{22, 80}, ScannedAt: now.Add(-time.Minute)})
	store.Append(Entry{Host: "h", Ports: []int{22, 80}, ScannedAt: now})

	deltas := Compare(store)
	if len(deltas) != 1 {
		t.Fatalf("expected 1 delta, got %d", len(deltas))
	}
	d := deltas[0]
	if len(d.Opened) != 0 || len(d.Closed) != 0 {
		t.Errorf("expected no opened/closed, got opened=%v closed=%v", d.Opened, d.Closed)
	}
}

func TestCompare_InsufficientEntries(t *testing.T) {
	store := &Store{}
	store.Append(Entry{Host: "lonely", Ports: []int{22}, ScannedAt: time.Now()})
	deltas := Compare(store)
	if len(deltas) != 0 {
		t.Errorf("expected no deltas for host with single entry, got %d", len(deltas))
	}
}

func TestCompare_EmptyStore(t *testing.T) {
	store := &Store{}
	deltas := Compare(store)
	if len(deltas) != 0 {
		t.Errorf("expected no deltas for empty store, got %d", len(deltas))
	}
}
