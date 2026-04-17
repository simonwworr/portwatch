package notify

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/user/portwatch/internal/state"
)

// EmailConfig holds SMTP configuration.
type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	To       []string
}

type emailChannel struct {
	cfg EmailConfig
}

// NewEmailChannel returns a Channel that sends notifications via SMTP.
func NewEmailChannel(cfg EmailConfig) Channel {
	return &emailChannel{cfg: cfg}
}

func (e *emailChannel) Send(host string, diff state.Diff) error {
	if len(e.cfg.To) == 0 {
		return fmt.Errorf("email: no recipients configured")
	}

	subject := fmt.Sprintf("[portwatch] Port changes detected on %s", host)
	body := formatBody(host, diff)

	msg := strings.Join([]string{
		"From: " + e.cfg.From,
		"To: " + strings.Join(e.cfg.To, ", "),
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=utf-8",
		"",
		body,
	}, "\r\n")

	addr := fmt.Sprintf("%s:%d", e.cfg.Host, e.cfg.Port)
	auth := smtp.PlainAuth("", e.cfg.Username, e.cfg.Password, e.cfg.Host)

	return smtp.SendMail(addr, auth, e.cfg.From, e.cfg.To, []byte(msg))
}
