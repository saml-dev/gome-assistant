// Package metricsui provides the read-only, embedded metrics console.
//
// It deliberately exposes HTML only. The package's small input boundary keeps
// the UI independent of the App package while allowing the App to supply its
// frozen automation catalog and current runtime state.
package metricsui

import (
	"errors"
	"net/http"
	"time"

	"saml.dev/gome-assistant/internal/metrics"
)

var errReaderUnavailable = errors.New("metrics reader is unavailable")

// Source is the committed metrics boundary consumed by the UI. Recorder
// satisfies this interface. Change notifications are invalidations only: the
// UI always re-reads committed state through Reader.
type Source interface {
	Reader() *metrics.Reader
	SubscribeChanges() (<-chan struct{}, func())
}

// Snapshot is the small live-state DTO shown in the Activity summary. It is
// intentionally separate from persisted telemetry, which remains browsable
// while no session is active.
type Snapshot struct {
	MetricsHealth  metrics.Health
	Connection     string
	SessionID      string
	SessionStarted time.Time
	SessionActive  bool
	SessionReady   bool
	LatestError    *metrics.Error
	Running        int
	Waiting        int
}

// Automation is one current, frozen logical automation. Configuration is the
// current in-memory definition, not a historical snapshot.
type Automation struct {
	Key            string
	Type           string
	Source         string
	Callback       string
	TriggerSummary string
	NextDue        time.Time
	Configuration  map[string]any
}

// Options supplies the two live App-owned values that metrics storage does not
// contain. The functions may be called for each request and must return
// snapshots safe for concurrent reading.
type Options struct {
	Source      Source
	Snapshot    func() Snapshot
	Automations func() []Automation
	Now         func() time.Time
}

// New returns a self-contained read-only HTTP handler. A nil Source produces
// a clear 503 response from data-bearing UI routes rather than panicking;
// normal App wiring always passes its metrics Recorder.
func New(options Options) http.Handler {
	if options.Now == nil {
		options.Now = time.Now
	}
	if options.Snapshot == nil {
		options.Snapshot = func() Snapshot { return Snapshot{} }
	}
	if options.Automations == nil {
		options.Automations = func() []Automation { return nil }
	}
	return &ui{options: options}
}

type ui struct {
	options Options
}

func (ui *ui) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writer.Header().Set("Allow", http.MethodGet)
		http.Error(writer, "metrics UI is read-only", http.StatusMethodNotAllowed)
		return
	}

	switch request.URL.Path {
	case "/", "/activity":
		if !ui.sourceAvailable() {
			ui.writeUnavailable(writer)
			return
		}
		ui.handleActivity(writer, request)
	case "/activity/fragment":
		if !ui.sourceAvailable() {
			ui.writeUnavailable(writer)
			return
		}
		ui.handleActivityFragment(writer, request)
	case "/activity/older":
		if !ui.sourceAvailable() {
			ui.writeUnavailable(writer)
			return
		}
		ui.handleOlder(writer, request)
	case "/activity/trace":
		if !ui.sourceAvailable() {
			ui.writeUnavailable(writer)
			return
		}
		ui.handleActivityTrace(writer, request)
	case "/automations":
		if !ui.sourceAvailable() {
			ui.writeUnavailable(writer)
			return
		}
		ui.handleAutomations(writer, request)
	case "/automations/fragment":
		if !ui.sourceAvailable() {
			ui.writeUnavailable(writer)
			return
		}
		ui.handleAutomationsFragment(writer, request)
	case "/events":
		if !ui.sourceAvailable() {
			ui.writeUnavailable(writer)
			return
		}
		ui.handleEvents(writer, request)
	case "/assets/htmx-2.0.10.min.js":
		ui.handleAsset(writer, request)
	default:
		http.NotFound(writer, request)
	}
}

func (ui *ui) sourceAvailable() bool {
	return ui.options.Source != nil && ui.options.Source.Reader() != nil
}

func (ui *ui) writeUnavailable(writer http.ResponseWriter) {
	http.Error(writer, errReaderUnavailable.Error(), http.StatusServiceUnavailable)
}
