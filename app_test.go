package gomeassistant

import (
	"context"
	stdhttp "net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewAppDoesNotPerformNetworkIO(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(stdhttp.HandlerFunc(func(stdhttp.ResponseWriter, *stdhttp.Request) {
		requests.Add(1)
	}))
	defer server.Close()

	app, err := NewApp(context.Background(), NewAppRequest{
		URL:         server.URL,
		HAAuthToken: "token",
	})
	require.NoError(t, err)
	require.NotNil(t, app)
	require.Nil(t, app.conn)
	require.Equal(t, int32(0), requests.Load())
	require.Equal(t, server.URL, app.baseURL.String())
	require.Equal(t, "token", app.authToken)
	require.Equal(t, "zone.home", app.homeZoneEntityID)
}

func TestNewAppRejectsInvalidArguments(t *testing.T) {
	tests := []struct {
		name    string
		request NewAppRequest
	}{
		{
			name: "missing URL",
			request: NewAppRequest{
				HAAuthToken: "token",
			},
		},
		{
			name: "missing token",
			request: NewAppRequest{
				URL: "http://example.com",
			},
		},
		{
			name: "malformed URL",
			request: NewAppRequest{
				URL:         "://",
				HAAuthToken: "token",
			},
		},
		{
			name: "unsupported URL scheme",
			request: NewAppRequest{
				URL:         "ws://example.com",
				HAAuthToken: "token",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := NewApp(context.Background(), tt.request)
			require.ErrorIs(t, err, ErrInvalidArgs)
			require.Nil(t, app)
		})
	}
}
