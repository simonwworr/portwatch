package history

import (
	"testing"
)

func TestBookmark_AddAndForHost(t *testing.T) {
	dir := t.TempDir()
	bs, _ := NewBookmarkStore(dir)
	bs.Add("host-a", 80, "web")
	bs.Add("host-a", 443, "tls")
	bs.Add("host-b", 22, "ssh")

	result := bs.ForHost("host-a")
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
}

func TestBookmark_ForHost_NoEntries(t *testing.T) {
	dir := t.TempDir()
	bs, _ := NewBookmarkStore(dir)
	if got := bs.ForHost("ghost"); len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestBookmark_Remove(t *testing.T) {
	dir := t.TempDir()
	bs, _ := NewBookmarkStore(dir)
	bs.Add("host-a", 80, "web")
	bs.Add("host-a", 443, "tls")
	bs.Remove("host-a", 80)

	result := bs.ForHost("host-a")
	if len(result) != 1 || result[0].Port != 443 {
		t.Fatalf("unexpected result after remove: %+v", result)
	}
}

func TestBookmark_SaveLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	bs, _ := NewBookmarkStore(dir)
	bs.Add("host-a", 8080, "dev")
	if err := bs.Save(); err != nil {
		t.Fatal(err)
	}

	bs2, err := NewBookmarkStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	all := bs2.All()
	if len(all) != 1 || all[0].Port != 8080 || all[0].Label != "dev" {
		t.Fatalf("unexpected loaded bookmarks: %+v", all)
	}
}

func TestBookmark_Load_MissingFile(t *testing.T) {
	dir := t.TempDir()
	bs, err := NewBookmarkStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(bs.All()) != 0 {
		t.Fatal("expected empty store")
	}
}
