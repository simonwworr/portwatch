package history

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAppendAndForHost(t *testing.T) {
	l := &Log{}
	l.Append("host1", []int{80, 443})
	l.Append("host2", []int{22})
	l.Append("host1", []int{80})

	entries := l.ForHost("host1")
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries for host1, got %d", len(entries))
	}
	if entries[0].Ports[1] != 443 {
		t.Errorf("expected port 443 in first entry")
	}
}

func TestForHost_NoEntries(t *testing.T) {
	l := &Log{}
	l.Append("host1", []int{80})
	if got := l.ForHost("ghost"); len(got) != 0 {
		t.Errorf("expected no entries, got %d", len(got))
	}
}

func TestSaveLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	l := &Log{}
	l.Append("192.168.1.1", []int{22, 80})

	if err := Save(path, l); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(loaded.Entries))
	}
	if loaded.Entries[0].Host != "192.168.1.1" {
		t.Errorf("unexpected host: %s", loaded.Entries[0].Host)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	l, err := Load("/nonexistent/path/history.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(l.Entries) != 0 {
		t.Errorf("expected empty log")
	}
}

func TestLoad_CorruptFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")
	os.WriteFile(path, []byte("not json{"), 0644)
	_, err := Load(path)
	if err == nil {
		t.Error("expected error for corrupt file")
	}
}
