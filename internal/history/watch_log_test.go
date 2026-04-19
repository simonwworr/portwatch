package history

import (
	"os"
	"testing"
	"time"
)

func TestAppendEvent_AndLoad(t *testing.T) {
	dir := t.TempDir()
	ev := WatchEvent{
		Host:      "192.168.1.1",
		Opened:    []int{80, 443},
		Closed:    []int{22},
		ScannedAt: time.Now().UTC().Truncate(time.Second),
	}
	if err := AppendEvent(dir, ev); err != nil {
		t.Fatalf("AppendEvent: %v", err)
	}
	wl, err := LoadWatchLog(dir)
	if err != nil {
		t.Fatalf("LoadWatchLog: %v", err)
	}
	if len(wl.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(wl.Events))
	}
	got := wl.Events[0]
	if got.Host != ev.Host {
		t.Errorf("host: got %s, want %s", got.Host, ev.Host)
	}
	if len(got.Opened) != 2 || got.Opened[0] != 80 {
		t.Errorf("opened ports mismatch: %v", got.Opened)
	}
}

func TestLoadWatchLog_MissingFile(t *testing.T) {
	dir := t.TempDir()
	wl, err := LoadWatchLog(dir)
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(wl.Events) != 0 {
		t.Errorf("expected empty log")
	}
}

func TestEventsForHost(t *testing.T) {
	dir := t.TempDir()
	hosts := []string{"10.0.0.1", "10.0.0.2", "10.0.0.1"}
	for _, h := range hosts {
		_ = AppendEvent(dir, WatchEvent{Host: h, ScannedAt: time.Now()})
	}
	wl, _ := LoadWatchLog(dir)
	evs := wl.EventsForHost("10.0.0.1")
	if len(evs) != 2 {
		t.Errorf("expected 2 events for 10.0.0.1, got %d", len(evs))
	}
	evs2 := wl.EventsForHost("10.0.0.2")
	if len(evs2) != 1 {
		t.Errorf("expected 1 event for 10.0.0.2, got %d", len(evs2))
	}
}

func TestLoadWatchLog_CorruptFile(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(watchLogPath(dir), []byte("not-json"), 0o644)
	_, err := LoadWatchLog(dir)
	if err == nil {
		t.Error("expected error for corrupt file")
	}
}
