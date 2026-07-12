package gomeassistant

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func newFireEventLifecycleTestApp(t *testing.T) *App {
	t.Helper()

	app, err := NewApp(context.Background(), NewAppRequest{
		URL:         "http://127.0.0.1:1",
		HAAuthToken: "token",
	})
	require.NoError(t, err)
	return app
}

func TestFireEventRejectsInactiveAppLifecycles(t *testing.T) {
	t.Run("before Start", func(t *testing.T) {
		app := newFireEventLifecycleTestApp(t)

		require.ErrorIs(t, app.FireEvent("test_event", nil), ErrAppNotRunning)
	})

	t.Run("after Start returns", func(t *testing.T) {
		app := newFireEventLifecycleTestApp(t)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		require.ErrorIs(t, app.Start(ctx), context.Canceled)
		require.ErrorIs(t, app.FireEvent("test_event", nil), ErrAppNotRunning)
	})

	t.Run("after Cleanup", func(t *testing.T) {
		app := newFireEventLifecycleTestApp(t)
		app.Cleanup()

		require.ErrorIs(t, app.FireEvent("test_event", nil), ErrAppClosed)
	})
}
