package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"portwatch/internal/history"
)

func writeHotspotHistory(t *testing.T, dir string) {
	t.Helper()
	now := time.Now()
	ms := history.NewMemoryStore()
	for i := 0; i < 3; i++ {
		_ = ms.Append(history.Entry{Host: "alpha", Ports: []int{80, 443}, Time: now.Add(time.Duration(i) * time.Hour)})
		_ = ms.Append(history.Entry{Host: "beta", Ports: []int{80, 8080}, Time: now.Add(time.Duration(i) * time.Hour)})
	}
	if err := history.Save(dir, ms); err != nil {
		t.Fatalf("failed to write hotspot history: %v", err)
	}
}

func TestRunHotspot_TextOutput(t *testing.T) {
	dir := t.TempDir()
	writeHotspotHistory(t, dir)

	buf := &bytes.Buffer{}
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runHotspot(dir, 2, "text")
	w.Close()
	os.Stdout = old
	buf.ReadFrom(r)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !containsStr(out, "80") {
		t.Errorf("expected port 80 in output, got: %s", out)
	}
}

func TestRunHotspot_JSONOutput(t *testing.T) {
	dir := t.TempDir()
	writeHotspotHistory(t, dir)

	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	err := runHotspot(dir, 2, "json")
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var results []map[string]interface{}
	if err := json.NewDecoder(r).Decode(&results); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(results) == 0 {
		t.Error("expected at least one result")
	}
}

func TestRunHotspot_MissingHistory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nonexistent")
	buf := &bytes.Buffer{}
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runHotspot(dir, 2, "text")
	w.Close()
	os.Stdout = old
	buf.ReadFrom(r)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsStr(buf.String(), "no history") {
		t.Errorf("expected 'no history' message, got: %s", buf.String())
	}
}
