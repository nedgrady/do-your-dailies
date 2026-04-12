package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"do-your-dailies/server/internal/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()
	app := New(nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestHealthCheckContentType(t *testing.T) {
	t.Parallel()
	app := New(nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
}

func TestHealthCheckBody(t *testing.T) {
	t.Parallel()
	app := New(nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, "healthy", rr.Body.String())
}

func TestHealthCheckWritesExplicitStatusBeforeBody(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	writer := testhelpers.NewRecordingResponseWriter()

	healthCheck(writer, req)

	assert.False(t, writer.WriteBeforeWriteHeader)
}

func TestHealthCheckWriterStatusCode(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	writer := testhelpers.NewRecordingResponseWriter()

	healthCheck(writer, req)

	assert.Equal(t, http.StatusOK, writer.StatusCode)
}

func TestHealthCheckWriterContentType(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	writer := testhelpers.NewRecordingResponseWriter()

	healthCheck(writer, req)

	assert.Equal(t, "text/plain; charset=utf-8", writer.Header().Get("Content-Type"))
}

func TestHealthCheckWriterBody(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	writer := testhelpers.NewRecordingResponseWriter()

	healthCheck(writer, req)

	assert.Equal(t, "healthy", writer.Body.String())
}
