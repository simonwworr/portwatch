package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"portwatch/internal/history"
)

func TestLabelAdd_AndList(t *testing.T) {
	dir := t.TempDir()
	s := history.NewLabelStore(dir)
	if err := s.Add("192.168.1.1", 80, "http"); err != nil {
		t.Fatalf("Add failed: %v", err)
	}
	loaded, err := history.LoadLabels(dir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	results := loaded.ForHost("192.168.1.1")
	if len(results) != 1 {
		t.Fatalf("expected 1 label, got %d", len(results))
	}
	if results[0].Label != "http" {
		t.Errorf("expected 'http', got %q", results[0].Label)
	}
}

func TestLabelAdd_MultiplePorts(t *testing.T) {
	dir := t.TempDir()
	s := history.NewLabelStore(dir)
	_ = s.Add("host", 22, "ssh")
	_ = s.Add("host", 443, "https")
	_ = s.Add("host", 8080, "alt-http")

	loaded, _ := history.LoadLabels(dir)
	if len(loaded.ForHost("host")) != 3 {
		t.Errorf("expected 3 labels")
	}
}

func TestLabelFile_CreatedInDir(t *testing.T) {
	dir := t.TempDir()
	s := history.NewLabelStore(dir)
	_ = s.Add("host", 80, "web")

	path := filepath.Join(dir, "labels.json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected labels.json to exist")
	}
}
