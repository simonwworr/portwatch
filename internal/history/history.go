package history

import (
	"encoding/json"
	"os"
	"time"
)

// Entry records a snapshot of open ports for a host at a point in time.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Host      string    `json:"host"`
	Ports     []int     `json:"ports"`
}

// Log is an ordered list of scan entries.
type Log struct {
	Entries []Entry `json:"entries"`
}

// Append adds a new entry to the log.
func (l *Log) Append(host string, ports []int) {
	l.Entries = append(l.Entries, Entry{
		Timestamp: time.Now().UTC(),
		Host:      host,
		Ports:     ports,
	})
}

// ForHost returns all entries for a given host.
func (l *Log) ForHost(host string) []Entry {
	var out []Entry
	for _, e := range l.Entries {
		if e.Host == host {
			out = append(out, e)
		}
	}
	return out
}

// Load reads a history log from disk. Returns an empty Log if the file does not exist.
func Load(path string) (*Log, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Log{}, nil
	}
	if err != nil {
		return nil, err
	}
	var l Log
	if err := json.Unmarshal(data, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// Save writes the log to disk as JSON.
func Save(path string, l *Log) error {
	data, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
