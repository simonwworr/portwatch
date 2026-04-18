package history

import (
	"os"
	"testing"
)

func TestAddTag_AndForHost(t *testing.T) {
	ts := &TagStore{}
	ts.AddTag("host1", "baseline", "initial scan", []int{22, 80})
	ts.AddTag("host2", "pre-deploy", "", []int{443})
	ts.AddTag("host1", "post-deploy", "after release", []int{22, 80, 443})

	tags := ts.ForHost("host1")
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags for host1, got %d", len(tags))
	}
	if tags[0].Name != "baseline" {
		t.Errorf("expected first tag name 'baseline', got %q", tags[0].Name)
	}
}

func TestForHost_NoEntries(t *testing.T) {
	ts := &TagStore{}
	if got := ts.ForHost("ghost"); len(got) != 0 {
		t.Errorf("expected empty, got %d entries", len(got))
	}
}

func TestFindByName_Found(t *testing.T) {
	ts := &TagStore{}
	ts.AddTag("host1", "v1", "version one", []int{80})

	tag := ts.FindByName("host1", "v1")
	if tag == nil {
		t.Fatal("expected to find tag, got nil")
	}
	if tag.Note != "version one" {
		t.Errorf("unexpected note: %q", tag.Note)
	}
}

func TestFindByName_NotFound(t *testing.T) {
	ts := &TagStore{}
	ts.AddTag("host1", "v1", "", []int{80})

	if got := ts.FindByName("host1", "v2"); got != nil {
		t.Errorf("expected nil, got %+v", got)
	}
}

func TestSaveLoadTags_RoundTrip(t *testing.T) {
	dir := t.TempDir()

	ts := &TagStore{}
	ts.AddTag("host1", "release", "prod release", []int{22, 443})

	if err := SaveTags(dir, ts); err != nil {
		t.Fatalf("SaveTags: %v", err)
	}

	loaded, err := LoadTags(dir)
	if err != nil {
		t.Fatalf("LoadTags: %v", err)
	}
	if len(loaded.Tags) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(loaded.Tags))
	}
	if loaded.Tags[0].Name != "release" {
		t.Errorf("unexpected tag name: %q", loaded.Tags[0].Name)
	}
}

func TestLoadTags_MissingFile(t *testing.T) {
	dir := t.TempDir()
	os.Remove(dir + "/tags.json")

	ts, err := LoadTags(dir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(ts.Tags) != 0 {
		t.Errorf("expected empty store")
	}
}
