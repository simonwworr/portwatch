package history

import (
	"testing"
	"time"
)

func buildScoreStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	now := time.Now()
	events := []WatchEvent{
		{Host: "host-a", Timestamp: now.Add(-2 * time.Hour), Opened: []int{80, 443}, Closed: []int{}},
		{Host: "host-a", Timestamp: now.Add(-1 * time.Hour), Opened: []int{8080}, Closed: []int{80}},
		{Host: "host-b", Timestamp: now.Add(-30 * time.Minute), Opened: []int{22}, Closed: []int{}},
	}
	for _, ev := range events {
		if err := AppendEvent(dir, ev); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func TestScore_Basic(t *testing.T) {
	dir := buildScoreStore(t)
	scores, err := Score(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(scores) != 2 {
		t.Fatalf("expected 2 hosts, got %d", len(scores))
	}
	// host-a has more changes
	if scores[0].Host != "host-a" {
		t.Errorf("expected host-a first, got %s", scores[0].Host)
	}
}

func TestScore_Openings(t *testing.T) {
	dir := buildScoreStore(t)
	scores, err := Score(dir)
	if err != nil {
		t.Fatal(err)
	}
	var ha HostScore
	for _, s := range scores {
		if s.Host == "host-a" {
			ha = s
		}
	}
	if ha.Openings != 3 {
		t.Errorf("expected 3 openings for host-a, got %d", ha.Openings)
	}
	if ha.Closings != 1 {
		t.Errorf("expected 1 closing for host-a, got %d", ha.Closings)
	}
}

func TestScore_EmptyStore(t *testing.T) {
	dir := t.TempDir()
	scores, err := Score(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(scores) != 0 {
		t.Errorf("expected empty scores, got %d", len(scores))
	}
}
