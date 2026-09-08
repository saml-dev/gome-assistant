package metricsui

import (
	"bytes"
	"embed"
	"net/http"
	"time"
)

//go:embed assets/htmx-2.0.10.min.js assets/htmx-sse-2.2.4.js
var assets embed.FS

func (ui *ui) handleAsset(writer http.ResponseWriter, request *http.Request) {
	var name string
	switch request.URL.Path {
	case "/assets/htmx-2.0.10.min.js":
		name = "assets/htmx-2.0.10.min.js"
	case "/assets/htmx-sse-2.2.4.js":
		name = "assets/htmx-sse-2.2.4.js"
	default:
		http.NotFound(writer, request)
		return
	}

	data, err := assets.ReadFile(name)
	if err != nil {
		http.Error(writer, "embedded asset is unavailable", http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	writer.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	http.ServeContent(writer, request, name, time.Time{}, bytes.NewReader(data))
}
