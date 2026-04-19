package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"portwatch/internal/history"
)

func writeSpikeHistory(t *testing.T, dir string) {
	t.Helper()
	s := &history.Store{}
	base := time.Now().Add(-10 * 24 * time.Hour)
	for i := 0; i < 5; i++ {
		s.Append(history.Entry{
			Host:      "h1",
			Ports:     []int{80, 443},
			ScannedAt: base.Add(time.Duration(i) * 24 * time.Hour),
		})
	}
	s.Append(history.Entry{
		Host:      "h1",
		Ports:     []int{80, 443, 8080, 8443, 9000, 9001, 9002, 9003},
		ScannedAt: base.Add(6 * 24 * time.Hour),
	})
	if err := history.Save(dir, s); err != nil {
		t.Fatal(err)
	}
}

func TestRunSpike_TextOutput(t *testing.T) {
	dir := t.TempDir()
	writeSpikeHistory(t, dir)
	if err := runSpike("h1", dir, 2.0, "text"); err != nil {
		t.Fatal(err)
	}
}

func TestRunSpike_JSONOutput(t *testing.T) {
	dir := t.TempDir()
	writeSpikeHistory(t, dir)

	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	err := runSpike("h1", dir, 2.0, "json")
	w.Close()
	os.Stdout = old
	if err != nil {
		t.Fatal(err)
	}
	var out []map[string]interface{}
	if err := json.NewDecoder(r).Decode(&out); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(out) == 0 {
		t.Error("expected spike results in json")
	}
}

func TestRunSpike_MissingHistory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "noexist")
	if err := runSpike("h1", dir, 2.0, "text"); err != nil {
		t.Fatal(err)
	}
}
