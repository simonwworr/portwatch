package alert_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/state"
)

func TestWebhookAlerter_Notify(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	a := alert.NewWebhookAlerter(ts.URL)
	diff := state.Diff{Opened: []int{9000}, Closed: []int{22}}
	if err := a.Notify("myhost", diff); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["host"] != "myhost" {
		t.Errorf("expected host=myhost, got %v", received["host"])
	}
}

func TestWebhookAlerter_Notify_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	a := alert.NewWebhookAlerter(ts.URL)
	diff := state.Diff{Opened: []int{80}}
	if err := a.Notify("host", diff); err == nil {
		t.Error("expected error on 500 response")
	}
}
