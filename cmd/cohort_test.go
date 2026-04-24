package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/history"
)

func writeCohortHistory(t *testing.T, dir string) {
	t.Helper()
	base := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	store := history.NewMemoryStore()
	store.Append("alpha", history.Entry{Time: base, Ports: []int{80, 443}})
	store.Append("beta", history.Entry{Time: base.Add(time.Hour), Ports: []int{80, 443, 8080}})
	if err := history.Save(dir, store); err != nil {
		t.Fatalf("save history: %v", err)
	}
}

func TestRunCohort_TextOutput(t *testing.T) {
	dir := t.TempDir()
	writeCohortHistory(t, dir)

	out := captureOutput(t, func() {
		if err := runCohort(dir, 1, "text"); err != nil {
			t.Fatalf("runCohort: %v", err)
		}
	})
	if !strings.Contains(out, "2024-03") {
		t.Errorf("expected cohort id in output, got: %s", out)
	}
	if !strings.Contains(out, "alpha") || !strings.Contains(out, "beta") {
		t.Errorf("expected both hosts in output, got: %s", out)
	}
}

func TestRunCohort_JSONOutput(t *testing.T) {
	dir := t.TempDir()
	writeCohortHistory(t, dir)

	out := captureOutput(t, func() {
		if err := runCohort(dir, 1, "json"); err != nil {
			t.Fatalf("runCohort: %v", err)
		}
	})
	var entries []map[string]interface{}
	if err := json.Unmarshal([]byte(out), &entries); err != nil {
		t.Fatalf("invalid JSON: %v\noutput: %s", err, out)
	}
	if len(entries) == 0 {
		t.Error("expected at least one cohort in JSON output")
	}
}

func TestRunCohort_MissingHistory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nonexistent")
	out := captureOutput(t, func() {
		if err := runCohort(dir, 1, "text"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(out, "no cohorts") {
		t.Errorf("expected no-cohorts message, got: %s", out)
	}
}

func TestRunCohort_MinHostsFilters(t *testing.T) {
	dir := t.TempDir()
	writeCohortHistory(t, dir)
	// Add a solo host in a different month
	april := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)
	store, _ := history.Load(dir)
	store.Append("solo", history.Entry{Time: april, Ports: []int{22}})
	_ = history.Save(dir, store)
	_ = os.Setenv("PORTWATCH_TEST", "1") // no-op, just for coverage

	out := captureOutput(t, func() {
		if err := runCohort(dir, 2, "text"); err != nil {
			t.Fatalf("runCohort: %v", err)
		}
	})
	if strings.Contains(out, "2024-04") {
		t.Errorf("solo cohort should be filtered out, got: %s", out)
	}
}
