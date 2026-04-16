package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/state"
)

// Alerter sends notifications about port changes.
type Alerter interface {
	Notify(host string, diff state.Diff) error
}

// LogAlerter writes alerts to a writer (default: stdout).
type LogAlerter struct {
	Out io.Writer
}

// NewLogAlerter creates a LogAlerter writing to stdout.
func NewLogAlerter() *LogAlerter {
	return &LogAlerter{Out: os.Stdout}
}

// Notify prints opened/closed port changes for a host.
func (a *LogAlerter) Notify(host string, diff state.Diff) error {
	ts := time.Now().Format(time.RFC3339)
	for _, p := range diff.Opened {
		fmt.Fprintf(a.Out, "[%s] ALERT host=%s port=%d status=opened\n", ts, host, p)
	}
	for _, p := range diff.Closed {
		fmt.Fprintf(a.Out, "[%s] ALERT host=%s port=%d status=closed\n", ts, host, p)
	}
	return nil
}

// HasChanges returns true when a diff contains any opened or closed ports.
func HasChanges(d state.Diff) bool {
	return len(d.Opened) > 0 || len(d.Closed) > 0
}
