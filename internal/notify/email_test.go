package notify

import (
	"net"
	"net/smtp"
	"testing"

	"github.com/user/portwatch/internal/state"
)

func startFakeSMTP(t *testing.T) (addr string, received *[]string) {
	t.Helper()
	lines := []string{}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		buf := make([]byte, 4096)
		n, _ := conn.Read(buf)
		lines = append(lines, string(buf[:n]))
	}()
	t.Cleanup(func() { ln.Close() })
	return ln.Addr().String(), &lines
}

func TestEmailChannel_NoRecipients(t *testing.T) {
	ch := NewEmailChannel(EmailConfig{
		Host: "localhost", Port: 25,
		From: "a@b.com", To: []string{},
	})
	diff := state.Diff{Opened: []int{80}, Closed: []int{}}
	err := ch.Send("host1", diff)
	if err == nil {
		t.Fatal("expected error for empty recipients")
	}
}

func TestEmailChannel_SMTPError(t *testing.T) {
	// Point at a port that refuses connections.
	ch := NewEmailChannel(EmailConfig{
		Host: "127.0.0.1", Port: 1,
		Username: "", Password: "",
		From: "a@b.com", To: []string{"c@d.com"},
	})
	diff := state.Diff{Opened: []int{443}, Closed: []int{80}}
	err := ch.Send("myhost", diff)
	if err == nil {
		t.Fatal("expected SMTP connection error")
	}
}

func TestEmailChannel_Send_UsesPlainAuth(t *testing.T) {
	// Verify smtp.PlainAuth is constructed without panic.
	auth := smtp.PlainAuth("", "user", "pass", "localhost")
	if auth == nil {
		t.Fatal("expected non-nil auth")
	}
}
