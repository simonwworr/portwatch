package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
	"time"

	"portwatch/internal/history"
)

func writeScoreHistory(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	now := time.Now()
	events := []history.WatchEvent{
		{Host: "alpha", Timestamp: now.Add(-1 * time.Hour), Opened: []int{80, 443}, Closed: []int{}},
		{Host: "beta", Timestamp: now.Add(-30 * time.Minute), Opened: []int{22}, Closed: []int{80}},
	}
	for _, ev := range events {
		if err := history.AppendEvent(dir, ev); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func TestRunScore_TextOutput(t *testing.T) {
	dir := writeScoreHistory(t)
	_ = os.Stdout
	if err := runScore(dir, "text"); err != nil {
		t.Fatal(err)
	}
}

func TestRunScore_JSONOutput(t *testing.T) {
	dir := writeScoreHistory(t)
	if err := runScore(dir, "json"); err != nil {
		t.Fatal(err)
	}
}

func TestRunScore_JSONStructure(t *testing.T) {
	dir := writeScoreHistory(t)
	scores, err := history.Score(dir)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(scores)
	var out []map[string]interface{}
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	if len(out) == 0 {
		t.Fatal("expected non-empty score list")
	}
	if _, ok := out[0]["host"]; !ok {
		t.Error("expected 'host' field")
	}
}

func TestRunScore_MissingHistory(t *testing.T) {
	dir := t.TempDir()
	var buf bytes.Buffer
	_ = buf
	if err := runScore(dir, "text"); err != nil {
		t.Fatal("expected no error for empty dir:", err)
	}
}
