package docs

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var docsGlobalStateLock sync.Mutex

func lockDocsGlobals(t *testing.T) {
	t.Helper()
	docsGlobalStateLock.Lock()
	t.Cleanup(func() {
		docsGlobalStateLock.Unlock()
	})
}

type recordingResponseWriter struct {
	header                 http.Header
	statusCode             int
	headerWritten          bool
	writeBeforeWriteHeader bool
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

	return len(payload), nil
}

func TestOpenAPISpecWritesExplicitStatusBeforeBody(t *testing.T) {
	t.Parallel()
	lockDocsGlobals(t)
	handler := NewHandler()
	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	writer := newRecordingResponseWriter()

	handler.openAPISpec(writer, req)

	assert.False(t, writer.writeBeforeWriteHeader)
}

func TestOpenAPISpecWriterStatusCode(t *testing.T) {
	t.Parallel()
	lockDocsGlobals(t)
	handler := NewHandler()
	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	writer := newRecordingResponseWriter()

	handler.openAPISpec(writer, req)

	assert.Equal(t, http.StatusOK, writer.statusCode)
}

func TestOpenAPISpecWriterContentType(t *testing.T) {
	t.Parallel()
	lockDocsGlobals(t)
	handler := NewHandler()
	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	writer := newRecordingResponseWriter()

	handler.openAPISpec(writer, req)

	assert.Equal(t, "application/yaml", writer.Header().Get("Content-Type"))
}

func TestOpenAPISpecReturns500WhenSpecMissing(t *testing.T) {
	t.Parallel()
	lockDocsGlobals(t)

	handler := NewHandler()
	original := openAPISpecPathFunc
	openAPISpecPathFunc = func() string { return filepath.Join(t.TempDir(), "missing-openapi.yaml") }
	t.Cleanup(func() {
		openAPISpecPathFunc = original
	})

	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	rr := httptest.NewRecorder()

	handler.openAPISpec(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestOpenAPISpecPathFallsBackToPrimaryCandidate(t *testing.T) {
	t.Parallel()
	lockDocsGlobals(t)

	originalWorkingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	temporaryWorkingDirectory := t.TempDir()

	if err := os.Chdir(temporaryWorkingDirectory); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalWorkingDirectory)
	})

	assert.Equal(t, filepath.Clean(filepath.Join("openapi", "openapi.yaml")), openAPISpecPath())
}

func TestSwaggerUIWritesExplicitStatusBeforeBody(t *testing.T) {
	t.Parallel()
	lockDocsGlobals(t)
	handler := NewHandler()
	req := httptest.NewRequest(http.MethodGet, "/swagger/", nil)
	writer := newRecordingResponseWriter()

	handler.swaggerUI(writer, req)

	assert.False(t, writer.writeBeforeWriteHeader)
}

func TestSwaggerUIWriterStatusCode(t *testing.T) {
	t.Parallel()
	lockDocsGlobals(t)
	handler := NewHandler()
	req := httptest.NewRequest(http.MethodGet, "/swagger/", nil)
	writer := newRecordingResponseWriter()

	handler.swaggerUI(writer, req)

	assert.Equal(t, http.StatusOK, writer.statusCode)
}

func TestSwaggerUIWriterContentType(t *testing.T) {
	t.Parallel()
	lockDocsGlobals(t)
	handler := NewHandler()
	req := httptest.NewRequest(http.MethodGet, "/swagger/", nil)
	writer := newRecordingResponseWriter()

	handler.swaggerUI(writer, req)

	assert.Equal(t, "text/html; charset=utf-8", writer.Header().Get("Content-Type"))
}
