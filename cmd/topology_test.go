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

func writeTopologyHistory(t *testing.T, dir string) {
	t.Helper()
	now := time.Now()
	s := history.NewMemoryStore()
	_ = s.Append(history.Entry{Host: "alpha", Ports: []int{80, 443}, ScannedAt: now})
	_ = s.Append(history.Entry{Host: "beta", Ports: []int{443, 9090}, ScannedAt: now})
	if err := history.Save(dir, s); err != nil {
		t.Fatalf("save history: %v", err)
	}
}

func TestRunTopology_TextOutput(t *testing.T) {
	dir := t.TempDir()
	writeTopologyHistory(t, dir)

	out := &bytes.Buffer{}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runTopology(dir, "text")
	w.Close()
	os.Stdout = oldStdout
	_, _ = bytes.NewBuffer(nil), out
	buf := &bytes.Buffer{}
	_, _ = buf.ReadFrom(r)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("alpha")) {
		t.Errorf("expected 'alpha' in output, got:\n%s", buf)
	}
	if !bytes.Contains(buf.Bytes(), []byte("443")) {
		t.Errorf("expected shared port 443 in output")
	}
}

func TestRunTopology_JSONOutput(t *testing.T) {
	dir := t.TempDir()
	writeTopologyHistory(t, dir)

	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w
	err := runTopology(dir, "json")
	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	buf := &bytes.Buffer{}
	_, _ = buf.ReadFrom(r)

	var res struct {
		Nodes []struct{ Host string }
		Edges []struct{ HostA, HostB string }
	}
	if err := json.Unmarshal(buf.Bytes(), &res); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, buf)
	}
	if len(res.Nodes) != 2 {
		t.Errorf("expected 2 nodes, got %d", len(res.Nodes))
	}
	if len(res.Edges) != 1 {
		t.Errorf("expected 1 edge, got %d", len(res.Edges))
	}
}

func TestRunTopology_MissingHistory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nonexistent")
	err := runTopology(dir, "text")
	if err != nil {
		t.Errorf("expected nil error for missing dir, got %v", err)
	}
}
