package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/state"
)

func TestPrint_TextFormat_NoChanges(t *testing.T) {
	r := New(map[string]state.Diff{})
	var buf bytes.Buffer
	r.Print(&buf, FormatText)
	if !strings.Contains(buf.String(), "No changes detected") {
		t.Errorf("expected no-changes message, got: %s", buf.String())
	}
}

func TestPrint_TextFormat_WithChanges(t *testing.T) {
	diffs := map[string]state.Diff{
		"localhost": {Opened: []int{8080}, Closed: []int{22}},
	}
	r := New(diffs)
	var buf bytes.Buffer
	r.Print(&buf, FormatText)
	out := buf.String()
	if !strings.Contains(out, "localhost") {
		t.Errorf("expected host in output, got: %s", out)
	}
	if !strings.Contains(out, "8080") {
		t.Errorf("expected opened port in output, got: %s", out)
	}
	if !strings.Contains(out, "22") {
		t.Errorf("expected closed port in output, got: %s", out)
	}
}

func TestPrint_JSONFormat(t *testing.T) {
	diffs := map[string]state.Diff{
		"192.168.1.1": {Opened: []int{443}, Closed: []int{}},
	}
	r := New(diffs)
	var buf bytes.Buffer
	r.Print(&buf, FormatJSON)
	out := buf.String()
	if !strings.Contains(out, "generated_at") {
		t.Errorf("expected generated_at in JSON, got: %s", out)
	}
	if !strings.Contains(out, "443") {
		t.Errorf("expected port 443 in JSON, got: %s", out)
	}
	if !strings.Contains(out, "192.168.1.1") {
		t.Errorf("expected host in JSON, got: %s", out)
	}
}

func TestIntSliceJSON_Empty(t *testing.T) {
	if intSliceJSON(nil) != "[]" {
		t.Error("expected [] for nil slice")
	}
}

func TestIntSliceJSON_Values(t *testing.T) {
	result := intSliceJSON([]int{80, 443, 8080})
	if result != "[80,443,8080]" {
		t.Errorf("unexpected JSON: %s", result)
	}
}
