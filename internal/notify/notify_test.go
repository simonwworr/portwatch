package notify_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/state"
)

type failChannel struct{}

func (f *failChannel) Send(_, _ string) error {
	return errors.New("send failed")
}

func TestDispatch_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	ch := notify.NewStdoutChannelWriter(&buf)
	d := notify.NewDispatcher(ch)
	err := d.Dispatch("localhost", state.Diff{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty diff, got: %s", buf.String())
	}
}

func TestDispatch_WithChanges(t *testing.T) {
	var buf bytes.Buffer
	ch := notify.NewStdoutChannelWriter(&buf)
	d := notify.NewDispatcher(ch)
	diff := state.Diff{Opened: []int{80, 443}, Closed: []int{8080}}
	err := d.Dispatch("example.com", diff)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "example.com") {
		t.Errorf("expected host in output, got: %s", out)
	}
	if !strings.Contains(out, "80") {
		t.Errorf("expected opened port in output, got: %s", out)
	}
	if !strings.Contains(out, "8080") {
		t.Errorf("expected closed port in output, got: %s", out)
	}
}

func TestDispatch_ChannelError(t *testing.T) {
	d := notify.NewDispatcher(&failChannel{})
	diff := state.Diff{Opened: []int{22}}
	err := d.Dispatch("host", diff)
	if err == nil {
		t.Fatal("expected error from failing channel")
	}
	if !strings.Contains(err.Error(), "send failed") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestStdoutChannel_Send(t *testing.T) {
	var buf bytes.Buffer
	ch := notify.NewStdoutChannelWriter(&buf)
	err := ch.Send("subject line", "body text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "subject line") {
		t.Errorf("expected subject in output")
	}
	if !strings.Contains(out, "body text") {
		t.Errorf("expected body in output")
	}
}
