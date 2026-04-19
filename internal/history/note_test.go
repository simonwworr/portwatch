package history

import (
	"testing"
	"os"
)

func TestNote_AddAndForHost(t *testing.T) {
	dir := t.TempDir()
	ns, err := NewNoteStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := ns.Add("host1", 80, "http open"); err != nil {
		t.Fatal(err)
	}
	if err := ns.Add("host1", 443, "https open"); err != nil {
		t.Fatal(err)
	}
	notes := ns.ForHost("host1")
	if len(notes) != 2 {
		t.Fatalf("expected 2 notes, got %d", len(notes))
	}
}

func TestNote_ForHost_NoEntries(t *testing.T) {
	dir := t.TempDir()
	ns, _ := NewNoteStore(dir)
	if got := ns.ForHost("ghost"); len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestNote_ForPort(t *testing.T) {
	dir := t.TempDir()
	ns, _ := NewNoteStore(dir)
	_ = ns.Add("host1", 80, "note a")
	_ = ns.Add("host1", 80, "note b")
	_ = ns.Add("host1", 443, "note c")
	got := ns.ForPort("host1", 80)
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestNote_EmptyText(t *testing.T) {
	dir := t.TempDir()
	ns, _ := NewNoteStore(dir)
	if err := ns.Add("host1", 80, ""); err == nil {
		t.Fatal("expected error for empty text")
	}
}

func TestNote_SaveLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	ns, _ := NewNoteStore(dir)
	_ = ns.Add("host1", 22, "ssh open")
	ns2, err := NewNoteStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(ns2.All()) != 1 {
		t.Fatalf("expected 1 note after reload, got %d", len(ns2.All()))
	}
	if ns2.All()[0].Text != "ssh open" {
		t.Errorf("unexpected text: %s", ns2.All()[0].Text)
	}
}

func TestNote_Load_MissingFile(t *testing.T) {
	dir := t.TempDir()
	_ = os.Remove(notePath(dir))
	ns, err := NewNoteStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(ns.All()) != 0 {
		t.Fatal("expected empty store")
	}
}
