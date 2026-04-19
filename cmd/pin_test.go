package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"portwatch/internal/history"
)

func TestPinAdd_CreatesPin(t *testing.T) {
	dir := t.TempDir()
	pinDir = dir

	err := runPinAdd(nil, []string{"host-a", "80,443", "web"})
	if err != nil {
		t.Fatal(err)
	}

	s := history.NewPinStore(dir)
	if err := s.Load(); err != nil {
		t.Fatal(err)
	}
	pins := s.ForHost("host-a")
	if len(pins) != 1 {
		t.Fatalf("expected 1 pin, got %d", len(pins))
	}
	if pins[0].Note != "web" {
		t.Errorf("unexpected note: %s", pins[0].Note)
	}
}

func TestPinAdd_InvalidPort(t *testing.T) {
	dir := t.TempDir()
	pinDir = dir
	err := runPinAdd(nil, []string{"host-a", "abc"})
	if err == nil {
		t.Error("expected error for invalid port")
	}
}

func TestPinList_AllHosts(t *testing.T) {
	dir := t.TempDir()
	pinDir = dir
	_ = runPinAdd(nil, []string{"host-a", "22", "ssh"})
	_ = runPinAdd(nil, []string{"host-b", "80", "http"})

	// just ensure no error
	if err := runPinList(nil, []string{}); err != nil {
		t.Error(err)
	}
}

func TestPinRemove_RemovesEntry(t *testing.T) {
	dir := t.TempDir()
	pinDir = dir
	_ = runPinAdd(nil, []string{"host-a", "9090", "metrics"})

	if err := runPinRemove(nil, []string{"host-a", "9090"}); err != nil {
		t.Fatal(err)
	}

	s := history.NewPinStore(dir)
	_ = s.Load()
	if len(s.ForHost("host-a")) != 0 {
		t.Error("expected pin to be removed")
	}
}

func TestPinFile_CreatedInDir(t *testing.T) {
	dir := t.TempDir()
	pinDir = dir
	_ = runPinAdd(nil, []string{"host-x", "3000"})

	path := filepath.Join(dir, "pins.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "host-x") {
		t.Error("pins.json missing host-x")
	}
}

func TestPinRemove_NonExistentHost(t *testing.T) {
	dir := t.TempDir()
	pinDir = dir

	// Removing a pin for a host that was never added should not error.
	if err := runPinRemove(nil, []string{"ghost-host", "8080"}); err != nil {
		t.Errorf("unexpected error removing non-existent pin: %v", err)
	}
}
