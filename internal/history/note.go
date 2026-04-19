package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Note struct {
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type NoteStore struct {
	dir   string
	notes []Note
}

func notePath(dir string) string {
	return filepath.Join(dir, "notes.json")
}

func NewNoteStore(dir string) (*NoteStore, error) {
	ns := &NoteStore{dir: dir}
	data, err := os.ReadFile(notePath(dir))
	if os.IsNotExist(err) {
		return ns, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &ns.notes); err != nil {
		return nil, err
	}
	return ns, nil
}

func (ns *NoteStore) Add(host string, port int, text string) error {
	if text == "" {
		return fmt.Errorf("note text must not be empty")
	}
	ns.notes = append(ns.notes, Note{
		Host:      host,
		Port:      port,
		Text:      text,
		CreatedAt: time.Now().UTC(),
	})
	return ns.save()
}

func (ns *NoteStore) ForHost(host string) []Note {
	var out []Note
	for _, n := range ns.notes {
		if n.Host == host {
			out = append(out, n)
		}
	}
	return out
}

func (ns *NoteStore) ForPort(host string, port int) []Note {
	var out []Note
	for _, n := range ns.notes {
		if n.Host == host && n.Port == port {
			out = append(out, n)
		}
	}
	return out
}

func (ns *NoteStore) All() []Note {
	return ns.notes
}

func (ns *NoteStore) save() error {
	if err := os.MkdirAll(ns.dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(ns.notes, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(notePath(ns.dir), data, 0644)
}
