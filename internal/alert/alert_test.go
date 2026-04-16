package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/state"
)

func TestLogAlerter_Notify_Opened(t *testing.T) {
	var buf bytes.Buffer
	a := &alert.LogAlerter{Out: &buf}
	diff := state.Diff{Opened: []int{8080}, Closed: []int{}}
	if err := a.Notify("localhost", diff); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "port=8080") {
		t.Errorf("expected port=8080 in output, got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "status=opened") {
		t.Errorf("expected status=opened in output")
	}
}

func TestLogAlerter_Notify_Closed(t *testing.T) {
	var buf bytes.Buffer
	a := &alert.LogAlerter{Out: &buf}
	diff := state.Diff{Opened: []int{}, Closed: []int{443}}
	_ = a.Notify("example.com", diff)
	if !strings.Contains(buf.String(), "status=closed") {
		t.Errorf("expected status=closed in output")
	}
}

func TestHasChanges(t *testing.T) {
	if alert.HasChanges(state.Diff{}) {
		t.Error("empty diff should have no changes")
	}
	if !alert.HasChanges(state.Diff{Opened: []int{22}}) {
		t.Error("diff with opened ports should have changes")
	}
	if !alert.HasChanges(state.Diff{Closed: []int{80}}) {
		t.Error("diff with closed ports should have changes")
	}
}
