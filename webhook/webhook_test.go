package webhook_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/tgbotkit/client"
	"github.com/tgbotkit/runtime/webhook"
)

func TestWebhook_ServeHTTP(t *testing.T) {
	opts := webhook.NewOptions()
	handler, err := webhook.NewWebhook(opts)
	if err != nil {
		t.Fatalf("failed to create webhook: %v", err)
	}

	update := client.Update{
		UpdateId: 123,
	}
	body, _ := json.Marshal(update)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	received, ok := handler.Receive(ctx)
	if !ok {
		t.Fatalf("expected to receive update from webhook")
	}

	if received.UpdateId != 123 {
		t.Errorf("expected update id 123, got %d", received.UpdateId)
	}
}

func TestWebhook_ServeHTTP_InvalidMethod(t *testing.T) {
	opts := webhook.NewOptions()
	handler, _ := webhook.NewWebhook(opts)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}