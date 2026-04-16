package notify

import (
	"fmt"
	"io"
	"os"
)

// StdoutChannel writes notifications to an io.Writer (default: os.Stdout).
type StdoutChannel struct {
	w io.Writer
}

// NewStdoutChannel creates a StdoutChannel writing to stdout.
func NewStdoutChannel() *StdoutChannel {
	return &StdoutChannel{w: os.Stdout}
}

// NewStdoutChannelWriter creates a StdoutChannel writing to w.
func NewStdoutChannelWriter(w io.Writer) *StdoutChannel {
	return &StdoutChannel{w: w}
}

// Send prints the subject and body to the writer.
func (s *StdoutChannel) Send(subject, body string) error {
	_, err := fmt.Fprintf(s.w, "--- %s ---\n%s\n", subject, body)
	return err
}
