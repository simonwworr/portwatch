package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Bookmark struct {
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	Label     string    `json:"label"`
	CreatedAt time.Time `json:"created_at"`
}

type BookmarkStore struct {
	dir  string
	items []Bookmark
}

func bookmarkPath(dir string) string {
	return filepath.Join(dir, "bookmarks.json")
}

func NewBookmarkStore(dir string) (*BookmarkStore, error) {
	bs := &BookmarkStore{dir: dir}
	data, err := os.ReadFile(bookmarkPath(dir))
	if os.IsNotExist(err) {
		return bs, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &bs.items); err != nil {
		return nil, err
	}
	return bs, nil
}

func (bs *BookmarkStore) Add(host string, port int, label string) {
	bs.items = append(bs.items, Bookmark{
		Host:      host,
		Port:      port,
		Label:     label,
		CreatedAt: time.Now(),
	})
}

func (bs *BookmarkStore) ForHost(host string) []Bookmark {
	var out []Bookmark
	for _, b := range bs.items {
		if b.Host == host {
			out = append(out, b)
		}
	}
	return out
}

func (bs *BookmarkStore) All() []Bookmark {
	return bs.items
}

func (bs *BookmarkStore) Remove(host string, port int) {
	var out []Bookmark
	for _, b := range bs.items {
		if b.Host == host && b.Port == port {
			continue
		}
		out = append(out, b)
	}
	bs.items = out
}

func (bs *BookmarkStore) Save() error {
	if err := os.MkdirAll(bs.dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(bs.items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(bookmarkPath(bs.dir), data, 0644)
}
