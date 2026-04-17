package history

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func makeStore() *Store {
	s := &Store{}
	t1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	s.Append(Entry{Host: "host1", ScannedAt: t1, OpenPorts: []int{80, 443}})
	s.Append(Entry{Host: "host1", ScannedAt: t2, OpenPorts: []int{80}})
	s.Append(Entry{Host: "host2", ScannedAt: t1, OpenPorts: []int{22}})
	return s
}

func TestExport_CSV(t *testing.T) {
	s := makeStore()
	var buf bytes.Buffer
	err := Export(s, ExportOptions{Host: "host1", Format: FormatCSV}, &buf)
	if err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "host,scanned_at,open_ports") {
		t.Error("missing CSV header")
	}
	if !strings.Contains(out, "80;443") {
		t.Error("expected port list 80;443")
	}
	if strings.Contains(out, "host2") {
		t.Error("host2 should not appear in host1 export")
	}
}

func TestExport_JSON(t *testing.T) {
	s := makeStore()
	var buf bytes.Buffer
	err := Export(s, ExportOptions{Host: "host1", Format: FormatJSON}, &buf)
	if err != nil {
		t.Fatal(err)
	}
	var entries []Entry
	if err := json.Unmarshal(buf.Bytes(), &entries); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}

func TestExport_Since(t *testing.T) {
	s := makeStore()
	since := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	var buf bytes.Buffer
	err := Export(s, ExportOptions{Host: "host1", Since: since, Format: FormatCSV}, &buf)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// header + 1 entry after since
	if len(lines) != 2 {
		t.Errorf("expected 2 lines (header+1), got %d", len(lines))
	}
}

func TestExport_UnknownFormat(t *testing.T) {
	s := makeStore()
	var buf bytes.Buffer
	err := Export(s, ExportOptions{Host: "host1", Format: "xml"}, &buf)
	if err == nil {
		t.Error("expected error for unknown format")
	}
}
