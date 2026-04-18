package history

import (
	"testing"
	"time"
)

func TestAnnotate_AddAndForHost(t *testing.T) {
	store := NewAnnotationStore()
	now := time.Now()

	if err := store.Add("192.168.1.1", now, "baseline established"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := store.Add("192.168.1.1", now.Add(time.Hour), "port 8080 opened intentionally"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	anns := store.ForHost("192.168.1.1")
	if len(anns) != 2 {
		t.Fatalf("expected 2 annotations, got %d", len(anns))
	}
	if anns[0].Note != "baseline established" {
		t.Errorf("unexpected note: %s", anns[0].Note)
	}
}

func TestAnnotate_ForHost_NoEntries(t *testing.T) {
	store := NewAnnotationStore()
	anns := store.ForHost("10.0.0.1")
	if len(anns) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(anns))
	}
}

func TestAnnotate_EmptyHost(t *testing.T) {
	store := NewAnnotationStore()
	err := store.Add("", time.Now(), "some note")
	if err == nil {
		t.Error("expected error for empty host")
	}
}

func TestAnnotate_EmptyNote(t *testing.T) {
	store := NewAnnotationStore()
	err := store.Add("10.0.0.1", time.Now(), "")
	if err == nil {
		t.Error("expected error for empty note")
	}
}

func TestAnnotate_All(t *testing.T) {
	store := NewAnnotationStore()
	now := time.Now()
	_ = store.Add("host-a", now, "note 1")
	_ = store.Add("host-b", now, "note 2")
	_ = store.Add("host-a", now.Add(time.Minute), "note 3")

	all := store.All()
	if len(all) != 3 {
		t.Errorf("expected 3 annotations, got %d", len(all))
	}
}
