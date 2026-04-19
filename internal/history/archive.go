package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ArchiveEntry represents a compressed snapshot of history for a host.
type ArchiveEntry struct {
	Host      string    `json:"host"`
	ArchivedAt time.Time `json:"archived_at"`
	Entries   []Entry   `json:"entries"`
}

// Archive moves all history entries for a host older than maxAge into an
// archive file, removing them from the live store.
func Archive(dir string, store *Store, host string, maxAge time.Duration) (int, error) {
	entries := store.ForHost(host)
	if len(entries) == 0 {
		return 0, nil
	}

	cutoff := time.Now().Add(-maxAge)
	var toArchive []Entry
	var toKeep []Entry

	for _, e := range entries {
		if e.Timestamp.Before(cutoff) {
			toArchive = append(toArchive, e)
		} else {
			toKeep = append(toKeep, e)
		}
	}

	if len(toArchive) == 0 {
		return 0, nil
	}

	archive := ArchiveEntry{
		Host:       host,
		ArchivedAt: time.Now(),
		Entries:    toArchive,
	}

	if err := os.MkdirAll(filepath.Join(dir, "archive"), 0755); err != nil {
		return 0, err
	}

	filename := fmt.Sprintf("%s_%d.json", host, time.Now().UnixNano())
	path := filepath.Join(dir, "archive", filename)

	f, err := os.Create(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(archive); err != nil {
		return 0, err
	}

	store.data[host] = toKeep
	return len(toArchive), nil
}

// LoadArchives returns all archived entries for a given host from the archive dir.
func LoadArchives(dir, host string) ([]ArchiveEntry, error) {
	archiveDir := filepath.Join(dir, "archive")
	matches, err := filepath.Glob(filepath.Join(archiveDir, host+"_*.json"))
	if err != nil {
		return nil, err
	}

	var results []ArchiveEntry
	for _, m := range matches {
		f, err := os.Open(m)
		if err != nil {
			continue
		}
		var a ArchiveEntry
		if err := json.NewDecoder(f).Decode(&a); err == nil {
			results = append(results, a)
		}
		f.Close()
	}
	return results, nil
}
