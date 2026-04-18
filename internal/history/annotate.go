package history

import (
	"fmt"
	"time"
)

// Annotation holds a user-defined note attached to a specific scan entry.
type Annotation struct {
	Host      string    `json:"host"`
	Timestamp time.Time `json:"timestamp"`
	Note      string    `json:"note"`
}

// AnnotationStore maps host -> list of annotations.
type AnnotationStore struct {
	Entries map[string][]Annotation `json:"entries"`
}

// NewAnnotationStore returns an empty AnnotationStore.
func NewAnnotationStore() *AnnotationStore {
	return &AnnotationStore{Entries: make(map[string][]Annotation)}
}

// Add attaches a note to the closest scan entry for the given host and time.
func (a *AnnotationStore) Add(host string, ts time.Time, note string) error {
	if host == "" {
		return fmt.Errorf("host must not be empty")
	}
	if note == "" {
		return fmt.Errorf("note must not be empty")
	}
	a.Entries[host] = append(a.Entries[host], Annotation{
		Host:      host,
		Timestamp: ts,
		Note:      note,
	})
	return nil
}

// ForHost returns all annotations for the given host.
func (a *AnnotationStore) ForHost(host string) []Annotation {
	return a.Entries[host]
}

// All returns every annotation across all hosts in chronological order.
func (a *AnnotationStore) All() []Annotation {
	var out []Annotation
	for _, anns := range a.Entries {
		out = append(out, anns...)
	}
	return out
}
