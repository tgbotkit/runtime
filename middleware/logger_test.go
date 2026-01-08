package middleware

import (
	"context"
	"errors"
	"testing"

	"github.com/tgbotkit/runtime/eventemitter"
)

type mockLogger struct {
	debugfCalled bool
	errorfCalled bool
}

func (m *mockLogger) Errorf(_ string, _ ...interface{}) { m.errorfCalled = true }
func (m *mockLogger) Fatalf(_ string, _ ...interface{}) {}
func (m *mockLogger) Fatal(_ ...interface{})            {}
func (m *mockLogger) Infof(_ string, _ ...interface{})  {}
func (m *mockLogger) Info(_ ...interface{})             {}
func (m *mockLogger) Warnf(_ string, _ ...interface{})  {}
func (m *mockLogger) Debugf(_ string, _ ...interface{}) { m.debugfCalled = true }
func (m *mockLogger) Debug(_ ...interface{})            {}

func TestLoggerMiddleware(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		l := &mockLogger{}
		mw := Logger(l)

		next := eventemitter.ListenerFunc(func(_ context.Context, _ any) error {
			return nil
		})

		err := mw.Handle(next).Handle(context.Background(), "test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !l.debugfCalled {
			t.Error("expected Debugf to be called")
		}
		if l.errorfCalled {
			t.Error("expected Errorf not to be called")
		}
	})

	t.Run("error", func(t *testing.T) {
		l := &mockLogger{}
		mw := Logger(l)

		expectedErr := errors.New("test error")
		next := eventemitter.ListenerFunc(func(_ context.Context, _ any) error {
			return expectedErr
		})

		err := mw.Handle(next).Handle(context.Background(), "test")
		if err != expectedErr {
			t.Fatalf("expected error %v, got %v", expectedErr, err)
		}

		if !l.debugfCalled {
			t.Error("expected Debugf to be called")
		}
		if !l.errorfCalled {
			t.Error("expected Errorf to be called")
		}
	})
}
