package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	app := New(nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestHealthCheckContentType(t *testing.T) {
	app := New(nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
}

func TestHealthCheckBody(t *testing.T) {
	app := New(nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, "healthy", rr.Body.String())
}

func TestHealthCheckWritesExplicitStatusBeforeBody(t *testing.T) {
	app := New(nil)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	writer := newRecordingResponseWriter()

	app.healthCheck(writer, req)

	assert.False(t, writer.writeBeforeWriteHeader)
}

func TestHealthCheckWriterStatusCode(t *testing.T) {
	app := New(nil)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	writer := newRecordingResponseWriter()

	app.healthCheck(writer, req)

	assert.Equal(t, http.StatusOK, writer.statusCode)
}

func TestHealthCheckWriterContentType(t *testing.T) {
	app := New(nil)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	writer := newRecordingResponseWriter()

	app.healthCheck(writer, req)

	assert.Equal(t, "text/plain; charset=utf-8", writer.Header().Get("Content-Type"))
}

func TestHealthCheckWriterBody(t *testing.T) {
	app := New(nil)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	writer := newRecordingResponseWriter()

	app.healthCheck(writer, req)

	assert.Equal(t, "healthy", writer.body.String())
}
