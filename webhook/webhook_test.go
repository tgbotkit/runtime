package webhook_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tgbotkit/client"
	"github.com/tgbotkit/runtime/webhook"
)

func TestNew(t *testing.T) {
	opts := webhook.NewOptions()
	wh, err := webhook.New(opts)
	require.NoError(t, err)
	assert.NotNil(t, wh)
	assert.NotNil(t, wh.UpdateChan())
}

func TestWebhook_ServeHTTP_Success(t *testing.T) {
	opts := webhook.NewOptions()
	wh, err := webhook.New(opts)
	require.NoError(t, err)

	update := client.Update{UpdateId: 123}
	body, _ := json.Marshal(update)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	wh.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	select {
	case u := <-wh.UpdateChan():
		assert.Equal(t, 123, u.UpdateId)
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for update")
	}
}

func TestWebhook_ServeHTTP_SecretToken(t *testing.T) {
	token := "my-secret-token"
	opts := webhook.NewOptions(webhook.WithToken(token))
	wh, err := webhook.New(opts)
	require.NoError(t, err)

	update := client.Update{UpdateId: 456}
	body, _ := json.Marshal(update)

	t.Run("Valid Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
		req.Header.Set(webhook.HeaderTelegramBotAPISecretToken, token)
		rr := httptest.NewRecorder()

		wh.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		
		select {
		case u := <-wh.UpdateChan():
			assert.Equal(t, 456, u.UpdateId)
		case <-time.After(time.Second):
			t.Fatal("timeout waiting for update")
		}
	})

	t.Run("Invalid Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
		req.Header.Set(webhook.HeaderTelegramBotAPISecretToken, "wrong-token")
		rr := httptest.NewRecorder()

		wh.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Missing Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		wh.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}

func TestWebhook_ServeHTTP_MethodNotAllowed(t *testing.T) {
	opts := webhook.NewOptions()
	wh, err := webhook.New(opts)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	wh.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestWebhook_ServeHTTP_BadRequest(t *testing.T) {
	opts := webhook.NewOptions()
	wh, err := webhook.New(opts)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("invalid-json")))
	rr := httptest.NewRecorder()

	wh.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestWebhook_ServeHTTP_Timeout(t *testing.T) {
	opts := webhook.NewOptions()
	wh, err := webhook.New(opts)
	require.NoError(t, err)

	update := client.Update{UpdateId: 789}
	body, _ := json.Marshal(update)

	// Fill the channel (buffer size 100)
	for i := 0; i < 100; i++ {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
		rr := httptest.NewRecorder()
		wh.ServeHTTP(rr, req)
		require.Equal(t, http.StatusOK, rr.Code)
	}

	// Channel is now full. Next request should block or timeout.
	
	// Create a context that is already canceled to simulate immediate timeout/client disconnect
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	wh.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusRequestTimeout, rr.Code)
}
