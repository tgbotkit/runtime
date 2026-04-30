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

func TestNew(t *testing.T) {
	opts := webhook.NewOptions()
	wh, err := webhook.New(opts)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}
	if wh == nil {
		t.Fatal("New() webhook is nil")
	}
	if wh.UpdateChan() == nil {
		t.Fatal("UpdateChan() is nil")
	}
}

func TestNew_InvalidBufferSize(t *testing.T) {
	opts := webhook.NewOptions(webhook.WithBufferSize(0))
	wh, err := webhook.New(opts)
	if err == nil {
		t.Fatal("New() error is nil, want validation error")
	}
	if wh != nil {
		t.Fatalf("New() webhook=%v, want nil", wh)
	}
}

func TestWebhook_ServeHTTP_Success(t *testing.T) {
	opts := webhook.NewOptions()
	wh, err := webhook.New(opts)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	update := client.Update{UpdateId: 123}
	body, err := json.Marshal(update)
	if err != nil {
		t.Fatalf("Marshal() unexpected error: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	wh.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("response code=%d, want %d", rr.Code, http.StatusOK)
	}

	select {
	case u := <-wh.UpdateChan():
		if u.UpdateId != 123 {
			t.Fatalf("update id=%d, want %d", u.UpdateId, 123)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for update")
	}
}

func TestWebhook_ServeHTTP_SecretToken(t *testing.T) {
	token := "my-secret-token"
	opts := webhook.NewOptions(webhook.WithToken(token))
	wh, err := webhook.New(opts)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	update := client.Update{UpdateId: 456}
	body, err := json.Marshal(update)
	if err != nil {
		t.Fatalf("Marshal() unexpected error: %v", err)
	}

	t.Run("valid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
		req.Header.Set(webhook.HeaderTelegramBotAPISecretToken, token)
		rr := httptest.NewRecorder()

		wh.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("response code=%d, want %d", rr.Code, http.StatusOK)
		}

		select {
		case u := <-wh.UpdateChan():
			if u.UpdateId != 456 {
				t.Fatalf("update id=%d, want %d", u.UpdateId, 456)
			}
		case <-time.After(time.Second):
			t.Fatal("timeout waiting for update")
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
		req.Header.Set(webhook.HeaderTelegramBotAPISecretToken, "wrong-token")
		rr := httptest.NewRecorder()

		wh.ServeHTTP(rr, req)
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("response code=%d, want %d", rr.Code, http.StatusUnauthorized)
		}
	})

	t.Run("missing token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		wh.ServeHTTP(rr, req)
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("response code=%d, want %d", rr.Code, http.StatusUnauthorized)
		}
	})
}

func TestWebhook_ServeHTTP_MethodNotAllowed(t *testing.T) {
	opts := webhook.NewOptions()
	wh, err := webhook.New(opts)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	wh.ServeHTTP(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("response code=%d, want %d", rr.Code, http.StatusMethodNotAllowed)
	}
}

func TestWebhook_ServeHTTP_BadRequest(t *testing.T) {
	opts := webhook.NewOptions()
	wh, err := webhook.New(opts)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("invalid-json")))
	rr := httptest.NewRecorder()

	wh.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("response code=%d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestWebhook_ServeHTTP_Timeout(t *testing.T) {
	opts := webhook.NewOptions()
	wh, err := webhook.New(opts)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	update := client.Update{UpdateId: 789}
	body, err := json.Marshal(update)
	if err != nil {
		t.Fatalf("Marshal() unexpected error: %v", err)
	}

	for i := 0; i < 100; i++ {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
		rr := httptest.NewRecorder()
		wh.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("fill request %d response code=%d, want %d", i, rr.Code, http.StatusOK)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	wh.ServeHTTP(rr, req)
	if rr.Code != http.StatusRequestTimeout {
		t.Fatalf("response code=%d, want %d", rr.Code, http.StatusRequestTimeout)
	}
}
