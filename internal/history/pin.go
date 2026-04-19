package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Pin struct {
	Host      string    `json:"host"`
	Ports     []int     `json:"ports"`
	Note      string    `json:"note"`
	PinnedAt  time.Time `json:"pinned_at"`
}

type PinStore struct {
	dir  string
	pins []Pin
}

func NewPinStore(dir string) *PinStore {
	return &PinStore{dir: dir}
}

func pinPath(dir string) string {
	return filepath.Join(dir, "pins.json")
}

func (s *PinStore) Load() error {
	data, err := os.ReadFile(pinPath(s.dir))
	if os.IsNotExist(err) {
		s.pins = nil
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.pins)
}

func (s *PinStore) Save() error {
	if err := os.MkdirAll(s.dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s.pins, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(pinPath(s.dir), data, 0644)
}

func (s *PinStore) Add(host string, ports []int, note string) {
	s.pins = append(s.pins, Pin{
		Host:     host,
		Ports:    ports,
		Note:     note,
		PinnedAt: time.Now().UTC(),
	})
}

func (s *PinStore) ForHost(host string) []Pin {
	var out []Pin
	for _, p := range s.pins {
		if p.Host == host {
			out = append(out, p)
		}
	}
	return out
}

func (s *PinStore) All() []Pin {
	return s.pins
}

func (s *PinStore) Remove(host string, port int) bool {
	var kept []Pin
	removed := false
	for _, p := range s.pins {
		if p.Host == host && containsInt(p.Ports, port) {
			removed = true
			continue
		}
		kept = append(kept, p)
	}
	s.pins = kept
	return removed
}

func containsInt(slice []int, v int) bool {
	for _, x := range slice {
		if x == v {
			return true
		}
	}
	return false
}
