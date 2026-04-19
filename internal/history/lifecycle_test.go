package history

import (
	"testing"
	"time"
)

func buildLifecycleStore() Store {
	now := time.Now().UTC().Truncate(time.Second)
	s := NewMemoryStore()
	s.Append("host1", Entry{Host: "host1", Ports: []int{80, 443}, Time: now.Add(-4 * time.Hour)})
	s.Append("host1", Entry{Host: "host1", Ports: []int{80, 443, 8080}, Time: now.Add(-3 * time.Hour)})
	s.Append("host1", Entry{Host: "host1", Ports: []int{80}, Time: now.Add(-2 * time.Hour)})
	s.Append("host1", Entry{Host: "host1", Ports: []int{80}, Time: now.Add(-1 * time.Hour)})
	return s
}

func TestLifecycle_OpenAndClosed(t *testing.T) {
	s := buildLifecycleStore()
	events := Lifecycle(s, "host1")

	closed := 0
	open := 0
	for _, e := range events {
		if e.ClosedAt != nil {
			closed++
		} else {
			open++
		}
	}
	if closed == 0 {
		t.Error("expected at least one closed lifecycle event")
	}
	if open == 0 {
		t.Error("expected at least one still-open lifecycle event")
	}
}

func TestLifecycle_Duration(t *testing.T) {
	s := buildLifecycleStore()
	events := Lifecycle(s, "host1")

	for _, e := range events {
		if e.ClosedAt != nil && e.Duration <= 0 {
			t.Errorf("port %d has non-positive duration %v", e.Port, e.Duration)
		}
	}
}

func TestLifecycle_EmptyStore(t *testing.T) {
	s := NewMemoryStore()
	events := Lifecycle(s, "host1")
	if len(events) != 0 {
		t.Errorf("expected 0 events, got %d", len(events))
	}
}

func TestLifecycle_Port8080_Closed(t *testing.T) {
	s := buildLifecycleStore()
	events := Lifecycle(s, "host1")

	for _, e := range events {
		if e.Port == 8080 && e.ClosedAt == nil {
			t.Error("port 8080 should be closed")
		}
	}
}
