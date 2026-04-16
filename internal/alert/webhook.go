package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/state"
)

// WebhookAlerter posts JSON payloads to a URL on port changes.
type WebhookAlerter struct {
	URL    string
	Client *http.Client
}

// NewWebhookAlerter creates a WebhookAlerter with a default HTTP client.
func NewWebhookAlerter(url string) *WebhookAlerter {
	return &WebhookAlerter{
		URL:    url,
		Client: &http.Client{Timeout: 10 * time.Second},
	}
}

type webhookPayload struct {
	Timestamp string `json:"timestamp"`
	Host      string `json:"host"`
	Opened    []int  `json:"opened"`
	Closed    []int  `json:"closed"`
}

// Notify sends a JSON POST request with the diff details.
func (w *WebhookAlerter) Notify(host string, diff state.Diff) error {
	payload := webhookPayload{
		Timestamp: time.Now().Format(time.RFC3339),
		Host:      host,
		Opened:    diff.Opened,
		Closed:    diff.Closed,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("alert: marshal payload: %w", err)
	}
	resp, err := w.Client.Post(w.URL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("alert: webhook post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("alert: webhook returned status %d", resp.StatusCode)
	}
	return nil
}
