package history

import (
	"testing"
	"time"
)

func buildWatchLog(events []WatchEvent) *WatchLog {
	return &WatchLog{Events: events}
}

func TestSummarizeWatchLog_Basic(t *testing.T) {
	log := buildWatchLog([]WatchEvent{
		{Host: "host1", Type: "opened", Ports: []int{80, 443}, Time: time.Now()},
		{Host: "host1", Type: "closed", Ports: []int{22}, Time: time.Now()},
		{Host: "host2", Type: "opened", Ports: []int{8080}, Time: time.Now()},
	})
	summaries := SummarizeWatchLog(log)
	if len(summaries) != 2 {
		t.Fatalf("expected 2 summaries, got %d", len(summaries))
	}
	if summaries[0].Host != "host1" {
		t.Errorf("expected host1 first, got %s", summaries[0].Host)
	}
	if summaries[0].Opened != 1 || summaries[0].Closed != 1 {
		t.Errorf("unexpected opened/closed counts: %+v", summaries[0])
	}
	if summaries[0].TotalEvents != 2 {
		t.Errorf("expected TotalEvents=2, got %d", summaries[0].TotalEvents)
	}
}

func TestSummarizeWatchLog_TopPorts(t *testing.T) {
	events := []WatchEvent{}
	for i := 0; i < 7; i++ {
		events = append(events, WatchEvent{
			Host:  "h",
			Type:  "opened",
			Ports: []int{80, 443, 22, 8080, 3306, 5432, 6379},
			Time:  time.Now(),
		})
	}
	log := buildWatchLog(events)
	summaries := SummarizeWatchLog(log)
	if len(summaries) != 1 {
		t.Fatalf("expected 1 summary")
	}
	if len(summaries[0].TopPorts) != 5 {
		t.Errorf("expected top 5 ports, got %d", len(summaries[0].TopPorts))
	}
}

func TestSummarizeWatchLog_Empty(t *testing.T) {
	log := buildWatchLog(nil)
	summaries := SummarizeWatchLog(log)
	if len(summaries) != 0 {
		t.Errorf("expected empty summaries")
	}
}

func TestSummarizeWatchLog_NoChange(t *testing.T) {
	log := buildWatchLog([]WatchEvent{
		{Host: "h", Type: "opened", Ports: []int{80}, Time: time.Now()},
		{Host: "h", Type: "opened", Ports: []int{80}, Time: time.Now()},
	})
	s := SummarizeWatchLog(log)
	if s[0].Opened != 2 || s[0].Closed != 0 {
		t.Errorf("unexpected counts: %+v", s[0])
	}
}
