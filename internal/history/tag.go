package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Tag associates a named label with a point-in-time snapshot for a host.
type Tag struct {
	Host      string    `json:"host"`
	Name      string    `json:"name"`
	Note      string    `json:"note,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	Ports     []int     `json:"ports"`
}

// TagStore holds all tags keyed by host.
type TagStore struct {
	Tags []Tag `json:"tags"`
}

// tagPath returns the file path for the tag store.
func tagPath(dir string) string {
	return filepath.Join(dir, "tags.json")
}

// LoadTags loads the tag store from dir, returning an empty store if missing.
func LoadTags(dir string) (*TagStore, error) {
	path := tagPath(dir)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &TagStore{}, nil
	}
	if err != nil {
		return nil, err
	}
	var ts TagStore
	if err := json.Unmarshal(data, &ts); err != nil {
		return nil, err
	}
	return &ts, nil
}

// SaveTags persists the tag store to dir.
func SaveTags(dir string, ts *TagStore) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(ts, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(tagPath(dir), data, 0o644)
}

// AddTag appends a tag to the store.
func (ts *TagStore) AddTag(host, name, note string, ports []int) {
	ts.Tags = append(ts.Tags, Tag{
		Host:      host,
		Name:      name,
		Note:      note,
		CreatedAt: time.Now().UTC(),
		Ports:     ports,
	})
}

// ForHost returns all tags for the given host.
func (ts *TagStore) ForHost(host string) []Tag {
	var out []Tag
	for _, t := range ts.Tags {
		if t.Host == host {
			out = append(out, t)
		}
	}
	return out
}

// FindByName returns the first tag matching host and name, or nil.
func (ts *TagStore) FindByName(host, name string) *Tag {
	for i := range ts.Tags {
		if ts.Tags[i].Host == host && ts.Tags[i].Name == name {
			return &ts.Tags[i]
		}
	}
	return nil
}
