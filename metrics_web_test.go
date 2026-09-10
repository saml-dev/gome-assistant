package gomeassistant

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"log/slog"
	"net"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"saml.dev/gome-assistant/internal/metrics"
)

func TestMetricsWebRuntimeServesActivityAndLogsActualAddress(t *testing.T) {
	app := newConfigurationTestApp(t)
	var logs bytes.Buffer
	previousLogger := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&logs, nil)))
	t.Cleanup(func() { slog.SetDefault(previousLogger) })

	runtime, err := newMetricsRuntime(app, filepath.Join(t.TempDir(), "metrics.db"), localWebUIAddress)
	require.NoError(t, err)
	t.Cleanup(func() { runtime.shutdown(nil) })
	require.NotEmpty(t, runtime.address)
	_, portText, err := net.SplitHostPort(runtime.address)
	require.NoError(t, err)
	port, err := strconv.Atoi(portText)
	require.NoError(t, err)
	require.GreaterOrEqual(t, port, defaultWebUIPort)

	client := &http.Client{Timeout: time.Second, Transport: &http.Transport{DisableCompression: true}}
	request, err := http.NewRequest(http.MethodGet, "http://"+runtime.address+"/", nil)
	require.NoError(t, err)
	request.Header.Set("Accept-Encoding", "gzip")
	response, err := client.Do(request)
	require.NoError(t, err)
	defer response.Body.Close()
	require.Equal(t, "gzip", response.Header.Get("Content-Encoding"))
	require.Contains(t, response.Header.Values("Vary"), "Accept-Encoding")
	compressed, err := gzip.NewReader(response.Body)
	require.NoError(t, err)
	defer compressed.Close()
	body, err := io.ReadAll(compressed)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)
	require.Contains(t, string(body), "Activity")
	require.Contains(t, logs.String(), runtime.address)
}

func TestMetricsWebRuntimeBindsAllIPv4InterfacesWhenExposed(t *testing.T) {
	app := newConfigurationTestApp(t)
	runtime, err := newMetricsRuntime(app, filepath.Join(t.TempDir(), "metrics.db"), exposedWebUIAddress)
	require.NoError(t, err)
	t.Cleanup(func() { runtime.shutdown(nil) })

	host, portText, err := net.SplitHostPort(runtime.address)
	require.NoError(t, err)
	require.Equal(t, "0.0.0.0", host)
	port, err := strconv.Atoi(portText)
	require.NoError(t, err)
	require.GreaterOrEqual(t, port, defaultWebUIPort)
}

func TestEnableWebUIStartRunsUntilCanceled(t *testing.T) {
	server, _, _ := newFakeHA(t, http.StatusOK)
	app := fakeApp(t, server.URL)
	require.NoError(t, app.EnableMetricsCollection(filepath.Join(t.TempDir(), "metrics.db")))
	require.NoError(t, app.EnableWebUI(false))

	runCtx, cancel := context.WithCancel(context.Background())
	runDone := make(chan error, 1)
	go func() { runDone <- app.Start(runCtx) }()
	t.Cleanup(func() {
		cancel()
		app.Cleanup()
	})

	deadline := time.Now().Add(time.Second)
	for !app.snapshotRuntimeStatus().SessionReady {
		select {
		case err := <-runDone:
			t.Fatalf("Start returned before cancellation: %v", err)
		default:
		}
		if time.Now().After(deadline) {
			t.Fatal("app did not become ready with the web UI enabled")
		}
		time.Sleep(5 * time.Millisecond)
	}

	address := app.snapshotRuntimeStatus().WebUIAddress
	require.NotEmpty(t, address)
	response, err := (&http.Client{Timeout: time.Second}).Get("http://" + address + "/")
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())
	require.Equal(t, http.StatusOK, response.StatusCode)
	select {
	case err := <-runDone:
		t.Fatalf("Start returned while its context remained active: %v", err)
	default:
	}

	cancel()
	require.ErrorIs(t, <-runDone, context.Canceled)
}

func TestMetricsWebRuntimeIncrementsPastBusyPort(t *testing.T) {
	occupied, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, occupied.Close()) })
	_, occupiedPortText, err := net.SplitHostPort(occupied.Addr().String())
	require.NoError(t, err)
	occupiedPort, err := strconv.Atoi(occupiedPortText)
	require.NoError(t, err)
	if occupiedPort == 65535 {
		t.Skip("kernel selected the final TCP port")
	}

	app := newConfigurationTestApp(t)
	runtime, err := newMetricsRuntime(app, filepath.Join(t.TempDir(), "metrics.db"), net.JoinHostPort("127.0.0.1", occupiedPortText))
	require.NoError(t, err)
	t.Cleanup(func() { runtime.shutdown(nil) })
	_, selectedPortText, err := net.SplitHostPort(runtime.address)
	require.NoError(t, err)
	selectedPort, err := strconv.Atoi(selectedPortText)
	require.NoError(t, err)
	require.Greater(t, selectedPort, occupiedPort)
}

func TestMetricsWebRuntimeExposedModeIncrementsPastBusyPort(t *testing.T) {
	occupied, err := net.Listen("tcp4", "0.0.0.0:0")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, occupied.Close()) })
	_, occupiedPortText, err := net.SplitHostPort(occupied.Addr().String())
	require.NoError(t, err)
	occupiedPort, err := strconv.Atoi(occupiedPortText)
	require.NoError(t, err)
	if occupiedPort == 65535 {
		t.Skip("kernel selected the final TCP port")
	}

	app := newConfigurationTestApp(t)
	runtime, err := newMetricsRuntime(app, filepath.Join(t.TempDir(), "metrics.db"), net.JoinHostPort("0.0.0.0", occupiedPortText))
	require.NoError(t, err)
	t.Cleanup(func() { runtime.shutdown(nil) })
	host, selectedPortText, err := net.SplitHostPort(runtime.address)
	require.NoError(t, err)
	require.Equal(t, "0.0.0.0", host)
	selectedPort, err := strconv.Atoi(selectedPortText)
	require.NoError(t, err)
	require.Greater(t, selectedPort, occupiedPort)
}

func TestMetricsWebRuntimeBindFailureClosesProvisionalRecorder(t *testing.T) {
	var logs bytes.Buffer
	previousLogger := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&logs, nil)))
	t.Cleanup(func() { slog.SetDefault(previousLogger) })

	databasePath := filepath.Join(t.TempDir(), "metrics.db")
	app := newConfigurationTestApp(t)
	require.NoError(t, app.EnableMetricsCollection(databasePath))
	require.NoError(t, app.EnableWebUI(false))
	app.webUIAddress = "not-an-address"
	err := app.Start(context.Background())
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "bind metrics web UI"))
	require.Contains(t, logs.String(), "metrics web UI failed to listen")
	require.Contains(t, logs.String(), "not-an-address")
	require.Nil(t, app.metricsRuntime)

	// Reopening confirms that the failed attempt released its writer handles.
	recorder, openErr := metrics.OpenRecorder(databasePath)
	require.NoError(t, openErr)
	require.NoError(t, recorder.Close(nil))
}
