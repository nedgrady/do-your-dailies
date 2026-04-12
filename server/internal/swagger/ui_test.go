package swagger

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"do-your-dailies/server/internal/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestSwaggerUIRouteReturns200(t *testing.T) {
	t.Parallel()
	router := newSwaggerRouter()

	req := httptest.NewRequest(http.MethodGet, "/swagger/", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestSwaggerUIRouteReferencesSpec(t *testing.T) {
	t.Parallel()
	router := newSwaggerRouter()

	req := httptest.NewRequest(http.MethodGet, "/swagger/", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Contains(t, rr.Body.String(), "/openapi.yaml")
}

func TestSwaggerUIRouteReturnsHTMLContentType(t *testing.T) {
	t.Parallel()
	router := newSwaggerRouter()

	req := httptest.NewRequest(http.MethodGet, "/swagger/", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, "text/html; charset=utf-8", rr.Header().Get("Content-Type"))
}

func TestSwaggerUIWritesExplicitStatusBeforeBody(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "/swagger/", nil)
	writer := testhelpers.NewRecordingResponseWriter()

	ui(writer, req)

	assert.False(t, writer.WriteBeforeWriteHeader)
}

func TestSwaggerUIWriterStatusCode(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "/swagger/", nil)
	writer := testhelpers.NewRecordingResponseWriter()

	ui(writer, req)

	assert.Equal(t, http.StatusOK, writer.StatusCode)
}

func TestSwaggerUIWriterContentType(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "/swagger/", nil)
	writer := testhelpers.NewRecordingResponseWriter()

	ui(writer, req)

	assert.Equal(t, "text/html; charset=utf-8", writer.Header().Get("Content-Type"))
}
