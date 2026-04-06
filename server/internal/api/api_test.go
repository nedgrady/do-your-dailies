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
	assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
	assert.Equal(t, "OK", rr.Body.String())
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
