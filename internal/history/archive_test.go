package history

import (
	"os"
	"testing"
	"time"
)

func buildArchiveStore(host string, ages []time.Duration) *Store {
	s := &Store{data: make(map[string][]Entry)}
	for _, age := range ages {
		s.data[host] = append(s.data[host], Entry{
			Host:      host,
			Timestamp: time.Now().Add(-age),
			Ports:     []int{80},
		})
	}
	return s
}

func TestArchive_MovesOldEntries(t *testing.T) {
	dir := t.TempDir()
	host := "10.0.0.1"
	s := buildArchiveStore(host, []time.Duration{
		72 * time.Hour,
		48 * time.Hour,
		1 * time.Hour,
	})

	n, err := Archive(dir, s, host, 24*time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 archived, got %d", n)
	}
	if len(s.ForHost(host)) != 1 {
		t.Errorf("expected 1 remaining entry, got %d", len(s.ForHost(host)))
	}
}

func TestArchive_NoOldEntries(t *testing.T) {
	dir := t.TempDir()
	host := "10.0.0.2"
	s := buildArchiveStore(host, []time.Duration{1 * time.Hour})

	n, err := Archive(dir, s, host, 24*time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 archived, got %d", n)
	}
}

func TestArchive_EmptyHost(t *testing.T) {
	dir := t.TempDir()
	s := &Store{data: make(map[string][]Entry)}
	n, err := Archive(dir, s, "ghost", time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0, got %d", n)
	}
}

func TestLoadArchives_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	host := "10.0.0.3"
	s := buildArchiveStore(host, []time.Duration{48 * time.Hour, 72 * time.Hour})

	_, err := Archive(dir, s, host, time.Hour)
	if err != nil {
		t.Fatalf("archive error: %v", err)
	}

	archives, err := LoadArchives(dir, host)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if len(archives) != 1 {
		t.Errorf("expected 1 archive file, got %d", len(archives))
	}
	if len(archives[0].Entries) != 2 {
		t.Errorf("expected 2 entries in archive, got %d", len(archives[0].Entries))
	}
}

func TestLoadArchives_MissingDir(t *testing.T) {
	dir := t.TempDir()
	archives, err := LoadArchives(dir, "nobody")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(archives) != 0 {
		t.Errorf("expected empty result")
	}
	_ = os.RemoveAll(dir)
}
