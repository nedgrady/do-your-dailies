package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenAPISpecRouteReturns200(t *testing.T) {
	app := New(nil)

	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	rr := httptest.NewRecorder()

	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestOpenAPISpecRouteReturnsSpecBody(t *testing.T) {
	app := New(nil)

	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	rr := httptest.NewRecorder()

	app.Router.ServeHTTP(rr, req)

	assert.Contains(t, rr.Body.String(), "openapi: 3.0.3")
}

func TestOpenAPISpecRouteReturnsYAMLContentType(t *testing.T) {
	app := New(nil)

	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	rr := httptest.NewRecorder()

	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, "application/yaml", rr.Header().Get("Content-Type"))
}
