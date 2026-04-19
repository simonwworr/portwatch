package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// WatchEvent records a single port-change event detected during a watch cycle.
type WatchEvent struct {
	Host    string    `json:"host"`
	Opened  []int     `json:"opened"`
	Closed  []int     `json:"closed"`
	ScannedAt time.Time `json:"scanned_at"`
}

// WatchLog is an ordered list of watch events persisted to disk.
type WatchLog struct {
	Events []WatchEvent `json:"events"`
}

func watchLogPath(dir string) string {
	return filepath.Join(dir, "watch_log.json")
}

// LoadWatchLog reads the watch log from dir, returning an empty log if missing.
func LoadWatchLog(dir string) (*WatchLog, error) {
	path := watchLogPath(dir)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &WatchLog{}, nil
	}
	if err != nil {
		return nil, err
	}
	var wl WatchLog
	if err := json.Unmarshal(data, &wl); err != nil {
		return nil, err
	}
	return &wl, nil
}

// AppendEvent adds a new event to the log and saves it.
func AppendEvent(dir string, ev WatchEvent) error {
	wl, err := LoadWatchLog(dir)
	if err != nil {
		return err
	}
	wl.Events = append(wl.Events, ev)
	data, err := json.MarshalIndent(wl, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(watchLogPath(dir), data, 0o644)
}

// EventsForHost returns all watch events for a specific host.
func (wl *WatchLog) EventsForHost(host string) []WatchEvent {
	var out []WatchEvent
	for _, ev := range wl.Events {
		if ev.Host == host {
			out = append(out, ev)
		}
	}
	return out
}
