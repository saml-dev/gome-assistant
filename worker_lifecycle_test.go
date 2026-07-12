package gomeassistant

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGoTrackedRejectsWorkAfterCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	app := &App{ctx: ctx}
	var invoked atomic.Int32

	cancel()
	app.goTracked(func() { invoked.Add(1) })

	require.Zero(t, invoked.Load())
}

func TestStartBindsSessionContextForWorkerAdmission(t *testing.T) {
	app, err := NewApp(context.Background(), NewAppRequest{
		URL:         "http://127.0.0.1:1",
		HAAuthToken: "token",
	})
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	require.ErrorIs(t, app.Start(ctx), context.Canceled)

	var invoked atomic.Int32
	app.goTracked(func() { invoked.Add(1) })
	require.Zero(t, invoked.Load())
}
