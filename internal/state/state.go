package state

import (
	"encoding/json"
	"os"
	"time"
)

// HostState holds the last known open ports for a host.
type HostState struct {
	Host      string    `json:"host"`
	Ports     []int     `json:"ports"`
	ScannedAt time.Time `json:"scanned_at"`
}

// StateStore maps host -> HostState.
type StateStore map[string]HostState

// Load reads state from a JSON file. Returns empty store if file doesn't exist.
func Load(path string) (StateStore, error) {
	store := make(StateStore)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return store, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}
	return store, nil
}

// Save writes the state store to a JSON file.
func Save(path string, store StateStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Diff compares previous and current ports, returning added and removed ports.
func Diff(previous, current []int) (added, removed []int) {
	prev := toSet(previous)
	curr := toSet(current)
	for p := range curr {
		if !prev[p] {
			added = append(added, p)
		}
	}
	for p := range prev {
		if !curr[p] {
			removed = append(removed, p)
		}
	}
	return
}

func toSet(ports []int) map[int]bool {
	s := make(map[int]bool, len(ports))
	for _, p := range ports {
		s[p] = true
	}
	return s
}
