package metricsui

import (
	"fmt"
	"net/http"
)

// handleEvents emits payload-free invalidations. A signal is sent immediately
// after connection so fresh and reconnected browsers fetch the current,
// committed view. Subsequent signals originate only from Source's post-commit
// broadcaster.
func (ui *ui) handleEvents(writer http.ResponseWriter, request *http.Request) {
	if ui.options.Source == nil {
		http.Error(writer, errReaderUnavailable.Error(), http.StatusServiceUnavailable)
		return
	}
	flusher, ok := writer.(http.Flusher)
	if !ok {
		http.Error(writer, "streaming is unavailable", http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-cache, no-transform")
	writer.Header().Set("Connection", "keep-alive")
	writer.WriteHeader(http.StatusOK)
	writeInvalidation(writer, flusher)

	changes, cancel := ui.options.Source.SubscribeChanges()
	defer cancel()
	for {
		select {
		case <-request.Context().Done():
			return
		case _, open := <-changes:
			if !open {
				return
			}
			writeInvalidation(writer, flusher)
		}
	}
}

func writeInvalidation(writer http.ResponseWriter, flusher http.Flusher) {
	// EventSource dispatches only after a data field, including an empty one.
	// The empty data value is deliberate: this is an invalidation, not a data
	// stream. HTMX fetches authoritative committed HTML after this event.
	_, _ = fmt.Fprint(writer, "event: change\ndata:\n\n")
	flusher.Flush()
}
