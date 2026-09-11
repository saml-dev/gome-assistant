package metricsui

import (
	"bytes"
	"embed"
	"net/http"
	"time"
)

//go:embed assets/htmx-2.0.10.min.js
var assets embed.FS

func (ui *ui) handleAsset(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/assets/htmx-2.0.10.min.js" {
		http.NotFound(writer, request)
		return
	}
	name := "assets/htmx-2.0.10.min.js"
	data, err := assets.ReadFile(name)
	if err != nil {
		http.Error(writer, "embedded asset is unavailable", http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	writer.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	http.ServeContent(writer, request, name, time.Time{}, bytes.NewReader(data))
}
