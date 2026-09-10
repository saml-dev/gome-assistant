package gomeassistant

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	stdhttp "net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	gorilla "github.com/gorilla/websocket"

	"saml.dev/gome-assistant/internal/metrics"
	"saml.dev/gome-assistant/internal/metricsui"
	"saml.dev/gome-assistant/websocket"
)

type metricsRuntime struct {
	recorder *metrics.Recorder
	server   *stdhttp.Server
	listener net.Listener
	address  string

	promoted atomic.Bool
	closing  atomic.Bool
	close    sync.Once
}

type metricsSession struct {
	runtime   *metricsRuntime
	recorder  *metrics.Recorder
	sessionID string
}

func newMetricsRuntime(app *App, databasePath, webAddress string) (*metricsRuntime, error) {
	recorder, err := metrics.OpenRecorder(databasePath)
	if err != nil {
		return nil, err
	}
	runtime := &metricsRuntime{recorder: recorder}
	if webAddress == "" {
		return runtime, nil
	}

	listener, err := listenMetricsWebUI(webAddress)
	if err != nil {
		slog.Error("metrics web UI failed to listen", "starting_address", webAddress, "error", err)
		_ = recorder.Close(nil)
		return nil, err
	}
	runtime.listener = listener
	runtime.address = listener.Addr().String()
	runtime.server = &stdhttp.Server{Handler: metricsui.WithCompression(metricsui.New(metricsui.Options{
		Source:      recorder,
		Snapshot:    app.metricsUISnapshot,
		Automations: app.metricsUIAutomations,
	}))}
	slog.Info("metrics web UI listening", "address", runtime.address)
	go func() {
		err := runtime.server.Serve(listener)
		if err != nil && !errors.Is(err, stdhttp.ErrServerClosed) && !runtime.closing.Load() {
			slog.Error("metrics web UI stopped unexpectedly", "error", err)
		}
	}()
	return runtime, nil
}

func listenMetricsWebUI(startAddress string) (net.Listener, error) {
	host, portText, err := net.SplitHostPort(startAddress)
	if err != nil {
		return nil, fmt.Errorf("bind metrics web UI: parse starting address %q: %w", startAddress, err)
	}
	startPort, err := strconv.Atoi(portText)
	if err != nil || startPort < 1 || startPort > 65535 {
		return nil, fmt.Errorf("bind metrics web UI: invalid starting port %q", portText)
	}

	var lastErr error
	network := "tcp"
	if host == "0.0.0.0" {
		network = "tcp4"
	}
	for port := startPort; port <= 65535; port++ {
		address := net.JoinHostPort(host, strconv.Itoa(port))
		listener, err := net.Listen(network, address)
		if err == nil {
			return listener, nil
		}
		if !errors.Is(err, syscall.EADDRINUSE) {
			return nil, fmt.Errorf("bind metrics web UI at %s: %w", address, err)
		}
		lastErr = err
	}
	return nil, fmt.Errorf("bind metrics web UI: no free port from %s through %s: %w",
		startAddress, net.JoinHostPort(host, "65535"), lastErr)
}

func (runtime *metricsRuntime) shutdown(final *metrics.Event) {
	if runtime == nil {
		return
	}
	runtime.close.Do(func() {
		runtime.closing.Store(true)
		if runtime.server != nil {
			_ = runtime.server.Close()
		}
		if err := runtime.recorder.Close(final); err != nil {
			slog.Error("metrics recorder shutdown incomplete", "error", err)
		}
	})
}

func (session *metricsSession) record(eventType metrics.EventType, payload any) bool {
	if session == nil || session.recorder == nil {
		return false
	}
	return session.recorder.Record(metrics.Event{
		Type:       eventType,
		OccurredAt: time.Now(),
		SessionID:  session.sessionID,
		Payload:    payload,
	})
}

func (app *App) beginMetricsSession() (*metricsSession, error) {
	if app.metricsDatabasePath == "" {
		return nil, nil
	}

	app.metricsMu.Lock()
	runtime := app.metricsRuntime
	if runtime == nil {
		var err error
		runtime, err = newMetricsRuntime(app, app.metricsDatabasePath, app.webUIAddress)
		if err != nil {
			app.metricsMu.Unlock()
			return nil, err
		}
		app.metricsRuntime = runtime
	}
	app.metricsMu.Unlock()

	if app.closed.Load() {
		app.closeMetricsRuntime(nil)
		return nil, ErrAppClosed
	}
	sessionID, err := metrics.NewID()
	if err != nil {
		if !runtime.promoted.Load() {
			app.clearAndCloseProvisional(runtime, nil)
		}
		return nil, err
	}
	app.statusMu.Lock()
	app.status.SessionID = sessionID
	app.status.SessionStartedAt = time.Now().UTC()
	app.status.WebUIAddress = runtime.address
	app.statusMu.Unlock()
	return &metricsSession{runtime: runtime, recorder: runtime.recorder, sessionID: sessionID}, nil
}

func (app *App) promoteMetricsRuntime(runtime *metricsRuntime) {
	if runtime != nil {
		runtime.promoted.Store(true)
	}
}

func (app *App) finishMetricsSession(session *metricsSession, ready, contextWasDone bool, runErr error) {
	if session == nil {
		return
	}
	reason := "runtime_error"
	switch {
	case !ready:
		reason = "startup_error"
	case contextWasDone:
		reason = "context_done"
	case isConnectionClosedError(runErr):
		reason = "connection_closed"
	}
	final := metrics.Event{
		Type:       metrics.EventSessionFinished,
		OccurredAt: time.Now(),
		SessionID:  session.sessionID,
		Payload: metrics.SessionFinishedPayload{
			Reason: reason,
			Error:  metrics.ErrorValue(runErr),
		},
	}

	if app.closed.Load() {
		app.closeMetricsRuntime(&final)
		return
	}
	if !session.runtime.promoted.Load() {
		app.clearAndCloseProvisional(session.runtime, &final)
		return
	}
	session.recorder.Record(final)
}

func isConnectionClosedError(err error) bool {
	if err == nil || errors.Is(err, ErrConnectionClosed) || errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
		return true
	}
	var closeError *gorilla.CloseError
	return errors.As(err, &closeError)
}

func (app *App) clearAndCloseProvisional(runtime *metricsRuntime, final *metrics.Event) {
	app.metricsMu.Lock()
	if app.metricsRuntime == runtime {
		app.metricsRuntime = nil
	}
	app.metricsMu.Unlock()
	runtime.shutdown(final)
}

func (app *App) closeMetricsRuntime(final *metrics.Event) {
	app.metricsMu.Lock()
	runtime := app.metricsRuntime
	app.metricsMu.Unlock()
	if runtime != nil {
		runtime.shutdown(final)
	}
}

type runtimeStatus struct {
	SessionID         string
	SessionStartedAt  time.Time
	SessionActive     bool
	SessionReady      bool
	Connection        string
	LatestError       *metrics.Error
	Running           int
	Waiting           int
	WebUIAddress      string
	PermanentlyClosed bool
}

type runtimeStatusSnapshot struct {
	MetricsEnabled    bool
	SessionID         string
	SessionStartedAt  time.Time
	SessionActive     bool
	SessionReady      bool
	Connection        string
	LatestError       *metrics.Error
	Running           int
	Waiting           int
	WebUIAddress      string
	PermanentlyClosed bool
	MetricsHealth     metrics.Health
}

func (app *App) snapshotRuntimeStatus() runtimeStatusSnapshot {
	app.statusMu.RLock()
	status := app.status
	app.statusMu.RUnlock()
	snapshot := runtimeStatusSnapshot{
		MetricsEnabled: app.metricsDatabasePath != "", SessionID: status.SessionID,
		SessionStartedAt: status.SessionStartedAt, SessionActive: status.SessionActive,
		SessionReady: status.SessionReady, Connection: status.Connection,
		Running: status.Running, Waiting: status.Waiting, WebUIAddress: status.WebUIAddress,
		PermanentlyClosed: status.PermanentlyClosed,
	}
	if status.LatestError != nil {
		value := *status.LatestError
		snapshot.LatestError = &value
	}
	app.metricsMu.Lock()
	runtime := app.metricsRuntime
	app.metricsMu.Unlock()
	if runtime != nil {
		snapshot.MetricsHealth = runtime.recorder.Health()
	}
	return snapshot
}

func (app *App) metricsUISnapshot() metricsui.Snapshot {
	status := app.snapshotRuntimeStatus()
	return metricsui.Snapshot{
		MetricsHealth:  status.MetricsHealth,
		Connection:     status.Connection,
		SessionID:      status.SessionID,
		SessionStarted: status.SessionStartedAt,
		SessionActive:  status.SessionActive,
		SessionReady:   status.SessionReady,
		LatestError:    status.LatestError,
		Running:        status.Running,
		Waiting:        status.Waiting,
	}
}

func (app *App) metricsUIAutomations() []metricsui.Automation {
	catalog := app.automationCatalogSnapshot()
	automations := make([]metricsui.Automation, len(catalog))
	for index, definition := range catalog {
		automations[index] = metricsui.Automation{
			Key:            definition.Key,
			Type:           definition.Type,
			Source:         definition.Source,
			Callback:       definition.Callback,
			TriggerSummary: definition.TriggerSummary,
			NextDue:        definition.NextDue,
			Configuration:  definition.Configuration,
		}
	}
	return automations
}

func (app *App) metricsReader() *metrics.Reader {
	app.metricsMu.Lock()
	defer app.metricsMu.Unlock()
	if app.metricsRuntime == nil {
		return nil
	}
	return app.metricsRuntime.recorder.Reader()
}

func (app *App) markSessionStarting() {
	app.statusMu.Lock()
	app.status.SessionActive = true
	app.status.SessionReady = false
	app.status.Connection = "connecting"
	app.status.LatestError = nil
	app.status.Running = 0
	app.status.Waiting = 0
	app.statusMu.Unlock()
}

func (app *App) markConnected() {
	app.statusMu.Lock()
	app.status.Connection = "connected"
	app.statusMu.Unlock()
}

func (app *App) markSessionReady() {
	app.statusMu.Lock()
	app.status.SessionReady = true
	app.statusMu.Unlock()
}

func (app *App) markSessionFinished(err error) {
	app.statusMu.Lock()
	app.status.SessionActive = false
	app.status.SessionReady = false
	app.status.Connection = "disconnected"
	app.status.Running = 0
	app.status.Waiting = 0
	app.status.LatestError = metrics.ErrorValue(err)
	app.statusMu.Unlock()
}

func (app *App) markPermanentlyClosed() {
	app.statusMu.Lock()
	app.status.PermanentlyClosed = true
	app.status.SessionActive = false
	app.status.Connection = "disconnected"
	app.statusMu.Unlock()
}

func (app *App) changeRunning(delta int) {
	app.statusMu.Lock()
	app.status.Running += delta
	app.statusMu.Unlock()
}

func (app *App) changeWaiting(delta int) {
	app.statusMu.Lock()
	app.status.Waiting += delta
	app.statusMu.Unlock()
}

func (app *App) setSessionContext(ctx context.Context) {
	app.lifecycleMu.Lock()
	app.ctx = ctx
	app.lifecycleMu.Unlock()
}

func (app *App) sessionContext() context.Context {
	app.lifecycleMu.RLock()
	defer app.lifecycleMu.RUnlock()
	return app.ctx
}

func (app *App) sessionStopping() bool {
	ctx := app.sessionContext()
	return ctx != nil && ctx.Err() != nil
}

func (app *App) setConnection(conn *websocket.Conn) {
	app.lifecycleMu.Lock()
	app.conn = conn
	app.lifecycleMu.Unlock()
}

func (app *App) connection() *websocket.Conn {
	app.lifecycleMu.RLock()
	defer app.lifecycleMu.RUnlock()
	return app.conn
}

func (app *App) beginWorkerAdmission() {
	app.workerAdmissionMu.Lock()
	app.workerAdmissionClosed = false
	app.workerAdmissionMu.Unlock()
}

func (app *App) endWorkerAdmission() {
	app.workerAdmissionMu.Lock()
	app.workerAdmissionClosed = true
	app.workerAdmissionMu.Unlock()
}
