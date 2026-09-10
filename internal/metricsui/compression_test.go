package metricsui

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAcceptsGzip(t *testing.T) {
	tests := []struct {
		name   string
		values []string
		want   bool
	}{
		{name: "absent"},
		{name: "gzip", values: []string{"gzip"}, want: true},
		{name: "case and whitespace", values: []string{"br, GZip ; Q = 0.4"}, want: true},
		{name: "multiple fields", values: []string{"br", "gzip"}, want: true},
		{name: "wildcard", values: []string{"br, *;q=0.2"}, want: true},
		{name: "explicit exclusion beats wildcard", values: []string{"gzip;q=0, *;q=1"}},
		{name: "zero wildcard", values: []string{"*;q=0"}},
		{name: "invalid quality", values: []string{"gzip;q=invalid"}},
		{name: "unsupported", values: []string{"br"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := acceptsGzip(test.values); got != test.want {
				t.Fatalf("acceptsGzip(%q) = %t, want %t", test.values, got, test.want)
			}
		})
	}
}

func TestAcceptsIdentity(t *testing.T) {
	tests := []struct {
		name   string
		values []string
		want   bool
	}{
		{name: "absent", want: true},
		{name: "implicit", values: []string{"gzip"}, want: true},
		{name: "explicit", values: []string{"identity;q=0.5"}, want: true},
		{name: "explicit exclusion", values: []string{"identity;q=0"}},
		{name: "wildcard exclusion", values: []string{"gzip;q=0, *;q=0"}},
		{name: "explicit beats wildcard", values: []string{"identity;q=1, *;q=0"}, want: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := acceptsIdentity(test.values); got != test.want {
				t.Fatalf("acceptsIdentity(%q) = %t, want %t", test.values, got, test.want)
			}
		})
	}
}

func TestActivityNegotiatesGzipCompression(t *testing.T) {
	_, handler := newTestUI(t, Snapshot{}, nil)
	handler = WithCompression(handler)
	plain := serve(handler, http.MethodGet, "/")
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Accept-Encoding", "br, gzip;q=0.8")
	compressed := httptest.NewRecorder()
	handler.ServeHTTP(compressed, request)

	if compressed.Code != http.StatusOK || compressed.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("compressed Activity response = %d, encoding %q", compressed.Code, compressed.Header().Get("Content-Encoding"))
	}
	if !headerContainsToken(compressed.Header(), "Vary", "Accept-Encoding") {
		t.Fatalf("compressed Activity Vary = %q", compressed.Header().Values("Vary"))
	}
	if compressed.Header().Get("Content-Length") != "" {
		t.Fatalf("compressed Activity Content-Length = %q", compressed.Header().Get("Content-Length"))
	}
	if compressed.Body.Len() >= plain.Body.Len() {
		t.Fatalf("compressed Activity size = %d, plain size = %d", compressed.Body.Len(), plain.Body.Len())
	}
	if got := gunzip(t, compressed.Body.Bytes()); !bytes.Equal(got, plain.Body.Bytes()) {
		t.Fatal("decompressed Activity response differs from identity response")
	}

	identityRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	identityRequest.Header.Set("Accept-Encoding", "gzip;q=0, *;q=1")
	identity := httptest.NewRecorder()
	handler.ServeHTTP(identity, identityRequest)
	if identity.Header().Get("Content-Encoding") != "" || !bytes.Equal(identity.Body.Bytes(), plain.Body.Bytes()) {
		t.Fatalf("explicitly excluded gzip response used encoding %q", identity.Header().Get("Content-Encoding"))
	}

	unacceptableRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	unacceptableRequest.Header.Set("Accept-Encoding", "gzip;q=0, identity;q=0")
	unacceptable := httptest.NewRecorder()
	handler.ServeHTTP(unacceptable, unacceptableRequest)
	if unacceptable.Code != http.StatusNotAcceptable || unacceptable.Header().Get("Content-Encoding") != "" {
		t.Fatalf("unacceptable encoding response = %d, encoding %q", unacceptable.Code, unacceptable.Header().Get("Content-Encoding"))
	}
}

func TestHeadNegotiatesGetHeadersWithoutWritingABody(t *testing.T) {
	_, handler := newTestUI(t, Snapshot{}, nil)
	handler = WithCompression(handler)
	request := httptest.NewRequest(http.MethodHead, "/", nil)
	request.Header.Set("Accept-Encoding", "gzip")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK || response.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("HEAD Activity response = %d, encoding %q", response.Code, response.Header().Get("Content-Encoding"))
	}
	if response.Body.Len() != 0 {
		t.Fatalf("HEAD Activity response wrote %d body bytes", response.Body.Len())
	}
	if !headerContainsToken(response.Header(), "Vary", "Accept-Encoding") {
		t.Fatalf("HEAD Activity Vary = %q", response.Header().Values("Vary"))
	}
}

func TestAssetCompressionPreservesContentAndRangeRequests(t *testing.T) {
	handler := WithCompression(New(Options{}))
	asset, err := assets.ReadFile("assets/htmx-2.0.10.min.js")
	if err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(http.MethodGet, "/assets/htmx-2.0.10.min.js", nil)
	request.Header.Set("Accept-Encoding", "gzip")
	compressed := httptest.NewRecorder()
	handler.ServeHTTP(compressed, request)
	if compressed.Code != http.StatusOK || compressed.Header().Get("Content-Encoding") != "gzip" || compressed.Header().Get("Content-Length") != "" {
		t.Fatalf("compressed asset = %d, encoding %q, length %q", compressed.Code, compressed.Header().Get("Content-Encoding"), compressed.Header().Get("Content-Length"))
	}
	if got := gunzip(t, compressed.Body.Bytes()); !bytes.Equal(got, asset) {
		t.Fatal("decompressed asset differs from embedded asset")
	}

	rangeRequest := httptest.NewRequest(http.MethodGet, "/assets/htmx-2.0.10.min.js", nil)
	rangeRequest.Header.Set("Accept-Encoding", "gzip")
	rangeRequest.Header.Set("Range", "bytes=0-9")
	ranged := httptest.NewRecorder()
	handler.ServeHTTP(ranged, rangeRequest)
	if ranged.Code != http.StatusPartialContent || ranged.Header().Get("Content-Encoding") != "" || !bytes.Equal(ranged.Body.Bytes(), asset[:10]) {
		t.Fatalf("range response = %d, encoding %q, body %q", ranged.Code, ranged.Header().Get("Content-Encoding"), ranged.Body.Bytes())
	}
}

func TestCompressionSkipsResponsesThatMustRemainUnchanged(t *testing.T) {
	tests := []struct {
		name            string
		status          int
		contentType     string
		contentEncoding string
		cacheControl    string
		contentRange    string
	}{
		{name: "reset content", status: http.StatusResetContent, contentType: "text/plain"},
		{name: "partial content", status: http.StatusPartialContent, contentType: "text/plain", contentRange: "bytes 0-3/4"},
		{name: "already encoded", status: http.StatusOK, contentType: "text/plain", contentEncoding: "br"},
		{name: "no transform", status: http.StatusOK, contentType: "text/plain", cacheControl: "private, no-transform"},
		{name: "binary", status: http.StatusOK, contentType: "application/octet-stream"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := WithCompression(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
				writer.Header().Set("Content-Type", test.contentType)
				if test.contentEncoding != "" {
					writer.Header().Set("Content-Encoding", test.contentEncoding)
				}
				if test.cacheControl != "" {
					writer.Header().Set("Cache-Control", test.cacheControl)
				}
				if test.contentRange != "" {
					writer.Header().Set("Content-Range", test.contentRange)
				}
				writer.WriteHeader(test.status)
				_, _ = writer.Write([]byte("body"))
			}))
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			request.Header.Set("Accept-Encoding", "gzip")
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if got := response.Header().Get("Content-Encoding"); got != test.contentEncoding {
				t.Fatalf("Content-Encoding = %q, want %q", got, test.contentEncoding)
			}
		})
	}
}

func gunzip(t *testing.T, body []byte) []byte {
	t.Helper()
	reader, err := gzip.NewReader(bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	decompressed, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	if err := reader.Close(); err != nil {
		t.Fatal(err)
	}
	return decompressed
}

func headerContainsToken(header http.Header, name, token string) bool {
	for _, field := range header.Values(name) {
		for _, value := range strings.Split(field, ",") {
			if strings.EqualFold(strings.TrimSpace(value), token) {
				return true
			}
		}
	}
	return false
}
