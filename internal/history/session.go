package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Session represents a contiguous watch run: start time, end time, and scans performed.
type Session struct {
	ID        string    `json:"id"`
	Host      string    `json:"host"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at"`
	ScanCount int       `json:"scan_count"`
	PortsSeen []int     `json:"ports_seen"`
}

func sessionPath(dir string) string {
	return filepath.Join(dir, "sessions.json")
}

// LoadSessions reads all sessions from dir. Returns empty slice if file missing.
func LoadSessions(dir string) ([]Session, error) {
	data, err := os.ReadFile(sessionPath(dir))
	if os.IsNotExist(err) {
		return []Session{}, nil
	}
	if err != nil {
		return nil, err
	}
	var sessions []Session
	if err := json.Unmarshal(data, &sessions); err != nil {
		return nil, err
	}
	return sessions, nil
}

// SaveSessions persists sessions to dir.
func SaveSessions(dir string, sessions []Session) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(sessions, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(sessionPath(dir), data, 0o644)
}

// AppendSession adds a session record and saves.
func AppendSession(dir string, s Session) error {
	sessions, err := LoadSessions(dir)
	if err != nil {
		return err
	}
	sessions = append(sessions, s)
	return SaveSessions(dir, sessions)
}

// SessionsForHost returns all sessions for a given host, sorted by StartedAt.
func SessionsForHost(sessions []Session, host string) []Session {
	var out []Session
	for _, s := range sessions {
		if s.Host == host {
			out = append(out, s)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].StartedAt.Before(out[j].StartedAt)
	})
	return out
}
