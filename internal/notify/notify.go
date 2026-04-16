package notify

import (
	"fmt"
	"strings"

	"github.com/user/portwatch/internal/state"
)

// Channel represents a notification delivery channel.
type Channel interface {
	Send(subject, body string) error
}

// Dispatcher fans out notifications to one or more channels.
type Dispatcher struct {
	channels []Channel
}

// NewDispatcher creates a Dispatcher with the given channels.
func NewDispatcher(channels ...Channel) *Dispatcher {
	return &Dispatcher{channels: channels}
}

// Dispatch formats and sends a notification for the given diff.
func (d *Dispatcher) Dispatch(host string, diff state.Diff) error {
	if len(diff.Opened) == 0 && len(diff.Closed) == 0 {
		return nil
	}

	subject := fmt.Sprintf("[portwatch] Port changes detected on %s", host)
	body := formatBody(host, diff)

	var errs []string
	for _, ch := range d.channels {
		if err := ch.Send(subject, body); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("notify dispatch errors: %s", strings.Join(errs, "; "))
	}
	return nil
}

func formatBody(host string, diff state.Diff) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Host: %s\n", host))
	if len(diff.Opened) > 0 {
		sb.WriteString(fmt.Sprintf("Opened ports: %v\n", diff.Opened))
	}
	if len(diff.Closed) > 0 {
		sb.WriteString(fmt.Sprintf("Closed ports: %v\n", diff.Closed))
	}
	return sb.String()
}
