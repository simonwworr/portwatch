package state

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadSave_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")

	store := StateStore{
		"localhost": {
			Host:      "localhost",
			Ports:     []int{80, 443},
			ScannedAt: time.Now().UTC().Truncate(time.Second),
		},
	}

	if err := Save(path, store); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	entry, ok := loaded["localhost"]
	if !ok {
		t.Fatal("expected localhost entry")
	}
	if len(entry.Ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(entry.Ports))
	}
}

func TestLoad_MissingFile(t *testing.T) {
	store, err := Load("/tmp/portwatch_nonexistent_state.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(store) != 0 {
		t.Errorf("expected empty store")
	}
	_ = os.Remove("/tmp/portwatch_nonexistent_state.json")
}

func TestDiff(t *testing.T) {
	added, removed := Diff([]int{80, 443}, []int{80, 8080})
	if len(added) != 1 || added[0] != 8080 {
		t.Errorf("expected added=[8080], got %v", added)
	}
	if len(removed) != 1 || removed[0] != 443 {
		t.Errorf("expected removed=[443], got %v", removed)
	}
}

func TestDiff_NoDifference(t *testing.T) {
	added, removed := Diff([]int{22, 80}, []int{22, 80})
	if len(added) != 0 || len(removed) != 0 {
		t.Errorf("expected no diff, got added=%v removed=%v", added, removed)
	}
}
