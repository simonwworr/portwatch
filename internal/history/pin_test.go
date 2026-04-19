package history

import (
	"os"
	"testing"
)

func TestPin_AddAndForHost(t *testing.T) {
	dir := t.TempDir()
	s := NewPinStore(dir)
	s.Add("host-a", []int{80, 443}, "web ports")
	s.Add("host-b", []int{22}, "ssh")

	pins := s.ForHost("host-a")
	if len(pins) != 1 {
		t.Fatalf("expected 1 pin, got %d", len(pins))
	}
	if pins[0].Note != "web ports" {
		t.Errorf("unexpected note: %s", pins[0].Note)
	}
}

func TestPin_ForHost_NoEntries(t *testing.T) {
	dir := t.TempDir()
	s := NewPinStore(dir)
	if got := s.ForHost("missing"); len(got) != 0 {
		t.Errorf("expected empty, got %v", got)
	}
}

func TestPin_SaveLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	s := NewPinStore(dir)
	s.Add("host-a", []int{8080}, "alt http")
	if err := s.Save(); err != nil {
		t.Fatal(err)
	}

	s2 := NewPinStore(dir)
	if err := s2.Load(); err != nil {
		t.Fatal(err)
	}
	pins := s2.ForHost("host-a")
	if len(pins) != 1 || pins[0].Ports[0] != 8080 {
		t.Errorf("round-trip failed: %+v", pins)
	}
}

func TestPin_Load_MissingFile(t *testing.T) {
	dir := t.TempDir()
	s := NewPinStore(dir)
	if err := s.Load(); err != nil {
		t.Errorf("expected nil on missing file, got %v", err)
	}
}

func TestPin_Remove(t *testing.T) {
	dir := t.TempDir()
	s := NewPinStore(dir)
	s.Add("host-a", []int{80}, "http")
	removed := s.Remove("host-a", 80)
	if !removed {
		t.Error("expected removed=true")
	}
	if len(s.ForHost("host-a")) != 0 {
		t.Error("expected no pins after remove")
	}
}

func TestPin_All(t *testing.T) {
	dir := t.TempDir()
	_ = os.MkdirAll(dir, 0755)
	s := NewPinStore(dir)
	s.Add("a", []int{1}, "")
	s.Add("b", []int{2}, "")
	if len(s.All()) != 2 {
		t.Errorf("expected 2 pins, got %d", len(s.All()))
	}
}
