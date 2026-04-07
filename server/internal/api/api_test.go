package api

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
)

type recordingResponseWriter struct {
	header                 http.Header
	statusCode             int
	headerWritten          bool
	writeBeforeWriteHeader bool
	body                   bytes.Buffer
}

func newRecordingResponseWriter() *recordingResponseWriter {
	return &recordingResponseWriter{header: make(http.Header)}
}

func (w *recordingResponseWriter) Header() http.Header {
	return w.header
}

func (w *recordingResponseWriter) WriteHeader(statusCode int) {
	w.headerWritten = true
	w.statusCode = statusCode
}

func (w *recordingResponseWriter) Write(payload []byte) (int, error) {
	if !w.headerWritten {
		w.writeBeforeWriteHeader = true
	}

	return w.body.Write(payload)
}

func TestHealthCheck(t *testing.T) {
	app := New(nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
	assert.Equal(t, "healthy", rr.Body.String())
}

func TestHealthCheckWritesExplicitStatusBeforeBody(t *testing.T) {
	app := New(nil)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	writer := newRecordingResponseWriter()

	app.healthCheck(writer, req)

	assert.True(t, writer.headerWritten)
	assert.False(t, writer.writeBeforeWriteHeader)
	assert.Equal(t, http.StatusOK, writer.statusCode)
	assert.Equal(t, "text/plain; charset=utf-8", writer.Header().Get("Content-Type"))
	assert.Equal(t, "healthy", writer.body.String())
}

func TestRouterLogsHealthRequests(t *testing.T) {
	originalLogger := middleware.DefaultLogger
	var logOutput bytes.Buffer
	middleware.DefaultLogger = middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger:  log.New(&logOutput, "", 0),
		NoColor: true,
	})
	t.Cleanup(func() {
		middleware.DefaultLogger = originalLogger
	})

	app := New(nil)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NotEmpty(t, logOutput.String())
}

func TestRouterRecoversFromPanics(t *testing.T) {
	app := New(nil)
	app.Router.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	rr := httptest.NewRecorder()

	assert.NotPanics(t, func() {
		app.Router.ServeHTTP(rr, req)
	})
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestListChores(t *testing.T) {
	app := New(nil)

	req := httptest.NewRequest(http.MethodGet, "/api/chores/", nil)
	rr := httptest.NewRecorder()

	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"message":"list chores"}`, rr.Body.String())
}

func TestCreateChore(t *testing.T) {
	app := New(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/chores/", nil)
	rr := httptest.NewRecorder()

	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"message":"create chore"}`, rr.Body.String())
}

func TestGetChore(t *testing.T) {
	app := New(nil)

	req := httptest.NewRequest(http.MethodGet, "/api/chores/123", nil)
	rr := httptest.NewRecorder()

	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"id":"123"}`, rr.Body.String())
}
