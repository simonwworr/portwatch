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

func writeShadowHistory(t *testing.T, dir string) {
	t.Helper()
	now := time.Now()
	s := history.NewMemoryStore()
	for i := 0; i < 5; i++ {
		s.Append(history.Entry{Host: "hostA", Time: now.Add(time.Duration(i) * time.Hour), Ports: []int{80, 443}})
	}
	s.Append(history.Entry{Host: "hostA", Time: now.Add(10 * time.Hour), Ports: []int{80, 9999}})
	if err := history.Save(dir, s); err != nil {
		t.Fatalf("save history: %v", err)
	}
}

func TestRunShadow_TextOutput(t *testing.T) {
	dir := t.TempDir()
	writeShadowHistory(t, dir)

	buf := &bytes.Buffer{}
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runShadow(dir, 3, 1, "text")
	w.Close()
	os.Stdout = old
	buf.ReadFrom(r)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !containsStr(out, "9999") {
		t.Errorf("expected port 9999 in output, got: %s", out)
	}
}

func TestRunShadow_JSONOutput(t *testing.T) {
	dir := t.TempDir()
	writeShadowHistory(t, dir)

	tmpOut := filepath.Join(dir, "out.json")
	f, _ := os.Create(tmpOut)
	old := os.Stdout
	os.Stdout = f
	err := runShadow(dir, 3, 1, "json")
	f.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(tmpOut)
	var results []history.ShadowPort
	if err := json.Unmarshal(data, &results); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(results) == 0 {
		t.Error("expected at least one shadow port in JSON output")
	}
}

func TestRunShadow_MissingHistory(t *testing.T) {
	dir := t.TempDir()
	err := runShadow(dir, 3, 1, "text")
	if err != nil {
		t.Fatalf("expected no error for missing history, got: %v", err)
	}
}
