package webhook_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/tgbotkit/client"
	"github.com/tgbotkit/runtime/webhook"
)

type mockClient struct {
	client.ClientWithResponsesInterface
	setWebhookFunc func(ctx context.Context, body client.SetWebhookJSONRequestBody) (*client.SetWebhookResponse, error)
}

func (m *mockClient) SetWebhookWithResponse(
	ctx context.Context,
	body client.SetWebhookJSONRequestBody,
	_ ...client.RequestEditorFn,
) (*client.SetWebhookResponse, error) {
	return m.setWebhookFunc(ctx, body)
}

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

func TestWebhook_ServeHTTP_TrailingData(t *testing.T) {
	opts := webhook.NewOptions()
	wh, err := webhook.New(opts)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"update_id":1} {}`))
	rr := httptest.NewRecorder()

	wh.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("response code=%d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestWebhook_ServeHTTP_BodyTooLarge(t *testing.T) {
	opts := webhook.NewOptions(webhook.WithMaxBodyBytes(4))
	wh, err := webhook.New(opts)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"update_id":1}`)))
	rr := httptest.NewRecorder()

	wh.ServeHTTP(rr, req)
	if rr.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("response code=%d, want %d", rr.Code, http.StatusRequestEntityTooLarge)
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

func TestWebhook_ServeHTTP_BufferFull(t *testing.T) {
	opts := webhook.NewOptions(webhook.WithBufferSize(1))
	wh, err := webhook.New(opts)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	update := client.Update{UpdateId: 789}
	body, err := json.Marshal(update)
	if err != nil {
		t.Fatalf("Marshal() unexpected error: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	wh.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("first response code=%d, want %d", rr.Code, http.StatusOK)
	}

	req = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	rr = httptest.NewRecorder()
	wh.ServeHTTP(rr, req)
	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("second response code=%d, want %d", rr.Code, http.StatusServiceUnavailable)
	}
}

func TestWebhook_SetWebhookOptions(t *testing.T) {
	allowedUpdates := []string{"message", "callback_query"}
	var gotBody client.SetWebhookJSONRequestBody
	tgClient := &mockClient{
		setWebhookFunc: func(_ context.Context, body client.SetWebhookJSONRequestBody) (*client.SetWebhookResponse, error) {
			gotBody = body

			return &client.SetWebhookResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &struct {
					Ok     client.SetWebhook200Ok `json:"ok"`
					Result bool                   `json:"result"`
				}{
					Ok:     true,
					Result: true,
				},
			}, nil
		},
	}

	wh, err := webhook.New(webhook.NewOptions(
		webhook.WithClient(tgClient),
		webhook.WithUrl("https://example.com/webhook"),
		webhook.WithToken("secret"),
		webhook.WithAllowedUpdates(allowedUpdates),
		webhook.WithDropPendingUpdates(true),
		webhook.WithMaxConnections(17),
	))
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	if err := wh.SetWebhook(context.Background()); err != nil {
		t.Fatalf("SetWebhook() unexpected error: %v", err)
	}

	if gotBody.Url != "https://example.com/webhook" {
		t.Fatalf("url=%q, want %q", gotBody.Url, "https://example.com/webhook")
	}
	if gotBody.SecretToken == nil || *gotBody.SecretToken != "secret" {
		t.Fatalf("secret token=%v, want secret", gotBody.SecretToken)
	}
	if gotBody.AllowedUpdates == nil {
		t.Fatal("allowed updates is nil")
	}
	if got, want := *gotBody.AllowedUpdates, allowedUpdates; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("allowed updates=%v, want %v", got, want)
	}
	if gotBody.DropPendingUpdates == nil || !*gotBody.DropPendingUpdates {
		t.Fatalf("drop pending updates=%v, want true", gotBody.DropPendingUpdates)
	}
	if gotBody.MaxConnections == nil || *gotBody.MaxConnections != 17 {
		t.Fatalf("max connections=%v, want 17", gotBody.MaxConnections)
	}
}

func TestWebhook_SetWebhookRegistrationMode(t *testing.T) {
	t.Run("explicitly disabled", func(t *testing.T) {
		wh, err := webhook.New(webhook.NewOptions(
			webhook.WithWebhookRegistrationEnabled(false),
			webhook.WithUrl(""),
		))
		if err != nil {
			t.Fatalf("New() unexpected error: %v", err)
		}

		if err := wh.Start(context.Background()); err != nil {
			t.Fatalf("Start() unexpected error: %v", err)
		}
	})

	t.Run("enabled without client", func(t *testing.T) {
		wh, err := webhook.New(webhook.NewOptions(webhook.WithUrl("https://example.com/webhook")))
		if err != nil {
			t.Fatalf("New() unexpected error: %v", err)
		}

		err = wh.Start(context.Background())
		if err == nil {
			t.Fatal("Start() error is nil, want missing client error")
		}
		if !strings.Contains(err.Error(), "client is required") {
			t.Fatalf("Start() error=%v, want missing client", err)
		}
	})

	t.Run("enabled without url", func(t *testing.T) {
		tgClient := &mockClient{
			setWebhookFunc: func(_ context.Context, _ client.SetWebhookJSONRequestBody) (*client.SetWebhookResponse, error) {
				return nil, errors.New("must not call SetWebhook")
			},
		}
		wh, err := webhook.New(webhook.NewOptions(webhook.WithClient(tgClient)))
		if err != nil {
			t.Fatalf("New() unexpected error: %v", err)
		}

		err = wh.Start(context.Background())
		if err == nil {
			t.Fatal("Start() error is nil, want missing url error")
		}
		if !strings.Contains(err.Error(), "url is required") {
			t.Fatalf("Start() error=%v, want missing url", err)
		}
	})
}
