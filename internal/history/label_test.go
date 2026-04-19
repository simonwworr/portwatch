package history

import (
	"os"
	"testing"
)

func TestAddLabel_AndForHost(t *testing.T) {
	dir := t.TempDir()
	s := NewLabelStore(dir)
	if err := s.Add("host1", 80, "http"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	results := s.ForHost("host1")
	if len(results) != 1 || results[0].Label != "http" {
		t.Errorf("expected label 'http', got %+v", results)
	}
}

func TestForHost_NoLabels(t *testing.T) {
	s := NewLabelStore(t.TempDir())
	if got := s.ForHost("ghost"); len(got) != 0 {
		t.Errorf("expected empty, got %v", got)
	}
}

func TestForPort_Labels(t *testing.T) {
	dir := t.TempDir()
	s := NewLabelStore(dir)
	_ = s.Add("host1", 80, "http")
	_ = s.Add("host1", 443, "https")
	results := s.ForPort("host1", 443)
	if len(results) != 1 || results[0].Label != "https" {
		t.Errorf("expected https label, got %+v", results)
	}
}

func TestSaveLoadLabels_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	s := NewLabelStore(dir)
	_ = s.Add("host1", 22, "ssh")
	_ = s.Add("host2", 3306, "mysql")

	loaded, err := LoadLabels(dir)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if len(loaded.All()) != 2 {
		t.Errorf("expected 2 labels, got %d", len(loaded.All()))
	}
}

func TestLoadLabels_MissingFile(t *testing.T) {
	dir := t.TempDir()
	s, err := LoadLabels(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.All()) != 0 {
		t.Errorf("expected empty store")
	}
}

func TestAddLabel_InvalidArgs(t *testing.T) {
	s := NewLabelStore(t.TempDir())
	if err := s.Add("", 80, "label"); err == nil {
		t.Error("expected error for empty host")
	}
	if err := s.Add("host", 0, "label"); err == nil {
		t.Error("expected error for zero port")
	}
	if err := s.Add("host", 80, ""); err == nil {
		t.Error("expected error for empty label")
	}
	_ = os.Remove(labelPath(s.dir))
}
