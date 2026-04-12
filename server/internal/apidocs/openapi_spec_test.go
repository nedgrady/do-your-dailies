package apidocs

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"do-your-dailies/server/internal/testhelpers"

	"github.com/go-chi/chi/v5"
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

func newSpecRouter() *chi.Mux {
	router := chi.NewRouter()
	RegisterRoutes(router)
	return router
}

func TestOpenAPISpecRouteReturns200(t *testing.T) {
	t.Parallel()
	router := newSpecRouter()

	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestOpenAPISpecRouteReturnsSpecBody(t *testing.T) {
	t.Parallel()
	router := newSpecRouter()

	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Contains(t, rr.Body.String(), "openapi: 3.0.3")
}

func TestOpenAPISpecRouteReturnsYAMLContentType(t *testing.T) {
	t.Parallel()
	router := newSpecRouter()

	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, "application/yaml", rr.Header().Get("Content-Type"))
}

func TestOpenAPISpecWritesExplicitStatusBeforeBody(t *testing.T) {
	t.Parallel()
	lockDocsGlobals(t)
	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	writer := testhelpers.NewRecordingResponseWriter()

	openAPISpec(writer, req)

	assert.False(t, writer.WriteBeforeWriteHeader)
}

func TestOpenAPISpecWriterStatusCode(t *testing.T) {
	t.Parallel()
	lockDocsGlobals(t)
	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	writer := testhelpers.NewRecordingResponseWriter()

	openAPISpec(writer, req)

	assert.Equal(t, http.StatusOK, writer.StatusCode)
}

func TestOpenAPISpecWriterContentType(t *testing.T) {
	t.Parallel()
	lockDocsGlobals(t)
	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	writer := testhelpers.NewRecordingResponseWriter()

	openAPISpec(writer, req)

	assert.Equal(t, "application/yaml", writer.Header().Get("Content-Type"))
}

func TestOpenAPISpecReturns500WhenSpecMissing(t *testing.T) {
	t.Parallel()
	lockDocsGlobals(t)

	original := openAPISpecPathFunc
	openAPISpecPathFunc = func() string { return filepath.Join(t.TempDir(), "missing-openapi.yaml") }
	t.Cleanup(func() {
		openAPISpecPathFunc = original
	})

	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	rr := httptest.NewRecorder()

	openAPISpec(rr, req)

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
