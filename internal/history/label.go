package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Label struct {
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	Label     string    `json:"label"`
	CreatedAt time.Time `json:"created_at"`
}

type LabelStore struct {
	dir    string
	labels []Label
}

func NewLabelStore(dir string) *LabelStore {
	return &LabelStore{dir: dir}
}

func (s *LabelStore) Add(host string, port int, label string) error {
	if host == "" || label == "" || port <= 0 {
		return fmt.Errorf("host, port, and label are required")
	}
	s.labels = append(s.labels, Label{
		Host:      host,
		Port:      port,
		Label:     label,
		CreatedAt: time.Now(),
	})
	return s.save()
}

func (s *LabelStore) ForHost(host string) []Label {
	var out []Label
	for _, l := range s.labels {
		if l.Host == host {
			out = append(out, l)
		}
	}
	return out
}

func (s *LabelStore) ForPort(host string, port int) []Label {
	var out []Label
	for _, l := range s.labels {
		if l.Host == host && l.Port == port {
			out = append(out, l)
		}
	}
	return out
}

func (s *LabelStore) All() []Label { return s.labels }

func labelPath(dir string) string {
	return filepath.Join(dir, "labels.json")
}

func (s *LabelStore) save() error {
	if err := os.MkdirAll(s.dir, 0755); err != nil {
		return err
	}
	f, err := os.Create(labelPath(s.dir))
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(s.labels)
}

func LoadLabels(dir string) (*LabelStore, error) {
	s := &LabelStore{dir: dir}
	f, err := os.Open(labelPath(dir))
	if os.IsNotExist(err) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return s, json.NewDecoder(f).Decode(&s.labels)
}
