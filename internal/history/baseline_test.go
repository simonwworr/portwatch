package history

import (
	"testing"
	"time"
)

func buildBaselineStore() *Store {
	s := &Store{}
	now := time.Now()
	s.Append("host-a", Entry{Host: "host-a", Ports: []int{80, 443}, ScannedAt: now.Add(-2 * time.Hour)})
	s.Append("host-a", Entry{Host: "host-a", Ports: []int{80, 443, 8080}, ScannedAt: now.Add(-1 * time.Hour)})
	s.Append("host-b", Entry{Host: "host-b", Ports: []int{22}, ScannedAt: now.Add(-1 * time.Hour)})
	return s
}

func TestBuildBaseline_LatestEntry(t *testing.T) {
	s := buildBaselineStore()
	bl := BuildBaseline(s)

	if len(bl["host-a"]) != 3 {
		t.Fatalf("expected 3 ports for host-a, got %d", len(bl["host-a"]))
	}
	if len(bl["host-b"]) != 1 || bl["host-b"][0] != 22 {
		t.Fatalf("expected [22] for host-b, got %v", bl["host-b"])
	}
}

func TestBuildBaseline_EmptyStore(t *testing.T) {
	s := &Store{}
	bl := BuildBaseline(s)
	if len(bl) != 0 {
		t.Fatalf("expected empty baseline, got %v", bl)
	}
}

func TestDetectDeviations_OpenedAndClosed(t *testing.T) {
	bl := Baseline{"host-a": {80, 443}}
	current := map[string][]int{"host-a": {80, 8080}}
	devs := DetectDeviations(bl, current, time.Now())

	if len(devs) != 1 {
		t.Fatalf("expected 1 deviation, got %d", len(devs))
	}
	d := devs[0]
	if len(d.Opened) != 1 || d.Opened[0] != 8080 {
		t.Errorf("expected opened=[8080], got %v", d.Opened)
	}
	if len(d.Closed) != 1 || d.Closed[0] != 443 {
		t.Errorf("expected closed=[443], got %v", d.Closed)
	}
}

func TestDetectDeviations_NoChange(t *testing.T) {
	bl := Baseline{"host-a": {80, 443}}
	current := map[string][]int{"host-a": {80, 443}}
	devs := DetectDeviations(bl, current, time.Now())
	if len(devs) != 0 {
		t.Fatalf("expected no deviations, got %d", len(devs))
	}
}

func TestDetectDeviations_NewHost(t *testing.T) {
	bl := Baseline{}
	current := map[string][]int{"host-new": {22, 80}}
	devs := DetectDeviations(bl, current, time.Now())
	if len(devs) != 1 {
		t.Fatalf("expected 1 deviation for new host, got %d", len(devs))
	}
	if len(devs[0].Opened) != 2 {
		t.Errorf("expected 2 opened ports, got %v", devs[0].Opened)
	}
}
