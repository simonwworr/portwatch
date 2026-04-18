package history

import (
	"testing"
	"time"
)

func buildSnapshotStore() *Store {
	s := &Store{}
	now := time.Now()
	s.Append("host-a", Entry{Ports: []int{80, 443}, Timestamp: now.Add(-2 * time.Hour)})
	s.Append("host-a", Entry{Ports: []int{80, 443, 8080}, Timestamp: now.Add(-1 * time.Hour)})
	s.Append("host-b", Entry{Ports: []int{22}, Timestamp: now})
	return s
}

func TestLatestSnapshot_ReturnsLast(t *testing.T) {
	s := buildSnapshotStore()
	snap, ok := LatestSnapshot(s, "host-a")
	if !ok {
		t.Fatal("expected snapshot")
	}
	if len(snap.Ports) != 3 {
		t.Errorf("expected 3 ports, got %d", len(snap.Ports))
	}
	if snap.Host != "host-a" {
		t.Errorf("unexpected host: %s", snap.Host)
	}
}

func TestLatestSnapshot_Missing(t *testing.T) {
	s := buildSnapshotStore()
	_, ok := LatestSnapshot(s, "ghost")
	if ok {
		t.Error("expected no snapshot for unknown host")
	}
}

func TestAllSnapshots_OncePerHost(t *testing.T) {
	s := buildSnapshotStore()
	snaps := AllSnapshots(s)
	if len(snaps) != 2 {
		t.Errorf("expected 2 snapshots, got %d", len(snaps))
	}
	hosts := map[string]bool{}
	for _, sn := range snaps {
		hosts[sn.Host] = true
	}
	if !hosts["host-a"] || !hosts["host-b"] {
		t.Error("missing expected hosts in snapshots")
	}
}
