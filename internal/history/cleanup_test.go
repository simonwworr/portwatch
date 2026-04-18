package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeFile(t *testing.T, dir, name string, age time.Duration) {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(`{}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if age > 0 {
		mt := time.Now().Add(-age)
		if err := os.Chtimes(p, mt, mt); err != nil {
			t.Fatal(err)
		}
	}
}

func TestCleanup_RemovesOldFiles(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "old.json", 48*time.Hour)
	writeFile(t, dir, "new.json", 0)

	res, err := Cleanup(CleanupOptions{Dir: dir, MaxAge: 24 * time.Hour})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Removed) != 1 {
		t.Fatalf("expected 1 removed, got %d", len(res.Removed))
	}
	if _, err := os.Stat(filepath.Join(dir, "old.json")); !os.IsNotExist(err) {
		t.Error("old.json should have been deleted")
	}
	if _, err := os.Stat(filepath.Join(dir, "new.json")); err != nil {
		t.Error("new.json should still exist")
	}
}

func TestCleanup_DryRun(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "old.json", 48*time.Hour)

	res, err := Cleanup(CleanupOptions{Dir: dir, MaxAge: 24 * time.Hour, DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Removed) != 1 {
		t.Fatalf("expected 1 in removed list, got %d", len(res.Removed))
	}
	if _, err := os.Stat(filepath.Join(dir, "old.json")); err != nil {
		t.Error("dry run should not delete files")
	}
}

func TestCleanup_IgnoresNonJSON(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "data.csv")
	_ = os.WriteFile(p, []byte("a,b"), 0o644)
	mt := time.Now().Add(-72 * time.Hour)
	_ = os.Chtimes(p, mt, mt)

	res, err := Cleanup(CleanupOptions{Dir: dir, MaxAge: 24 * time.Hour})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Removed) != 0 {
		t.Error("non-json files should not be removed")
	}
}

func TestCleanup_MissingDir(t *testing.T) {
	res, err := Cleanup(CleanupOptions{Dir: "/nonexistent/path", MaxAge: time.Hour})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Removed) != 0 {
		t.Error("expected no removals for missing dir")
	}
}
