package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSwaggerUIRouteReturns200(t *testing.T) {
	app := New(nil)

	req := httptest.NewRequest(http.MethodGet, "/swagger/", nil)
	rr := httptest.NewRecorder()

	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestSwaggerUIRouteReferencesSpec(t *testing.T) {
	app := New(nil)

	req := httptest.NewRequest(http.MethodGet, "/swagger/", nil)
	rr := httptest.NewRecorder()

	app.Router.ServeHTTP(rr, req)

	assert.Contains(t, rr.Body.String(), "/openapi.yaml")
}

func TestSwaggerUIRouteReturnsHTMLContentType(t *testing.T) {
	app := New(nil)

	req := httptest.NewRequest(http.MethodGet, "/swagger/", nil)
	rr := httptest.NewRecorder()

	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, "text/html; charset=utf-8", rr.Header().Get("Content-Type"))
}
