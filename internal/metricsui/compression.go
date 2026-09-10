package metricsui

import (
	"compress/gzip"
	"net/http"
	"strconv"
	"strings"
)

type compressionHandler struct {
	next http.Handler
}

// WithCompression negotiates gzip for finite UI and asset responses. The SSE
// invalidation stream is always left uncompressed so Flush remains immediate.
func WithCompression(next http.Handler) http.Handler {
	return compressionHandler{next: next}
}

func (handler compressionHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// EventSource messages must reach the browser as soon as each invalidation
	// is flushed. The stream is deliberately tiny, so compression has no value.
	if request.URL.Path == "/events" {
		handler.next.ServeHTTP(writer, request)
		return
	}
	if request.Method == http.MethodHead {
		writer = &headResponseWriter{ResponseWriter: writer}
		request = request.Clone(request.Context())
		request.Method = http.MethodGet
	}

	addVary(writer.Header(), "Accept-Encoding")
	fields := request.Header.Values("Accept-Encoding")
	gzipAccepted := acceptsGzip(fields)
	identityAccepted := acceptsIdentity(fields)
	if !gzipAccepted && !identityAccepted {
		http.Error(writer, "no acceptable response encoding", http.StatusNotAcceptable)
		return
	}
	if request.Header.Get("Range") != "" {
		if !identityAccepted {
			http.Error(writer, "range requests require identity encoding", http.StatusNotAcceptable)
			return
		}
		handler.next.ServeHTTP(writer, request)
		return
	}
	if !gzipAccepted {
		handler.next.ServeHTTP(writer, request)
		return
	}

	compressed := &gzipResponseWriter{
		ResponseWriter: writer,
		gzip:           gzip.NewWriter(writer),
	}
	defer compressed.close()
	handler.next.ServeHTTP(compressed, request)
}

type headResponseWriter struct {
	http.ResponseWriter
}

func (writer *headResponseWriter) Unwrap() http.ResponseWriter {
	return writer.ResponseWriter
}

func (writer *headResponseWriter) Write(body []byte) (int, error) {
	return len(body), nil
}

type gzipResponseWriter struct {
	http.ResponseWriter
	gzip        *gzip.Writer
	wroteHeader bool
	compress    bool
}

func (writer *gzipResponseWriter) Unwrap() http.ResponseWriter {
	return writer.ResponseWriter
}

func (writer *gzipResponseWriter) WriteHeader(status int) {
	if writer.wroteHeader {
		return
	}
	writer.wroteHeader = true
	writer.compress = statusAllowsBody(status) &&
		writer.Header().Get("Content-Encoding") == "" &&
		writer.Header().Get("Content-Range") == "" &&
		!headerContainsDirective(writer.Header(), "Cache-Control", "no-transform") &&
		compressibleContentType(writer.Header().Get("Content-Type"))
	if writer.compress {
		writer.Header().Del("Content-Length")
		writer.Header().Set("Content-Encoding", "gzip")
	}
	writer.ResponseWriter.WriteHeader(status)
}

func (writer *gzipResponseWriter) Write(body []byte) (int, error) {
	if !writer.wroteHeader {
		writer.WriteHeader(http.StatusOK)
	}
	if writer.compress {
		return writer.gzip.Write(body)
	}
	return writer.ResponseWriter.Write(body)
}

func (writer *gzipResponseWriter) close() {
	if writer.compress {
		_ = writer.gzip.Close()
	}
}

func statusAllowsBody(status int) bool {
	return status >= http.StatusOK &&
		status != http.StatusNoContent &&
		status != http.StatusResetContent &&
		status != http.StatusPartialContent &&
		status != http.StatusNotModified
}

func acceptsGzip(fields []string) bool {
	return acceptedQuality(fields, "gzip", 0) > 0
}

func acceptsIdentity(fields []string) bool {
	return acceptedQuality(fields, "identity", 1) > 0
}

func acceptedQuality(fields []string, wanted string, defaultQuality float64) float64 {
	explicit, explicitQuality := false, float64(0)
	wildcard, wildcardQuality := false, float64(0)
	for _, field := range fields {
		for _, item := range strings.Split(field, ",") {
			parts := strings.Split(item, ";")
			coding := strings.ToLower(strings.TrimSpace(parts[0]))
			if coding == "" {
				continue
			}
			quality := float64(1)
			for _, parameter := range parts[1:] {
				name, value, found := strings.Cut(parameter, "=")
				if !found || !strings.EqualFold(strings.TrimSpace(name), "q") {
					continue
				}
				parsed, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
				if err != nil || parsed < 0 || parsed > 1 {
					quality = 0
				} else {
					quality = parsed
				}
			}
			switch {
			case coding == wanted:
				explicit = true
				explicitQuality = max(explicitQuality, quality)
			case coding == "*":
				wildcard = true
				wildcardQuality = max(wildcardQuality, quality)
			}
		}
	}
	if explicit {
		return explicitQuality
	}
	if wildcard {
		return wildcardQuality
	}
	return defaultQuality
}

func addVary(header http.Header, value string) {
	for _, field := range header.Values("Vary") {
		for _, existing := range strings.Split(field, ",") {
			if strings.EqualFold(strings.TrimSpace(existing), value) {
				return
			}
		}
	}
	header.Add("Vary", value)
}

func compressibleContentType(value string) bool {
	mediaType := strings.ToLower(strings.TrimSpace(strings.SplitN(value, ";", 2)[0]))
	return strings.HasPrefix(mediaType, "text/") ||
		mediaType == "application/javascript" ||
		mediaType == "application/json" ||
		mediaType == "application/xml" ||
		strings.HasSuffix(mediaType, "+json") ||
		strings.HasSuffix(mediaType, "+xml")
}

func headerContainsDirective(header http.Header, name, directive string) bool {
	for _, field := range header.Values(name) {
		for _, value := range strings.Split(field, ",") {
			if strings.EqualFold(strings.TrimSpace(strings.SplitN(value, "=", 2)[0]), directive) {
				return true
			}
		}
	}
	return false
}
