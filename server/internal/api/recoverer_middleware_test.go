package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouterRecoversFromPanics(t *testing.T) {
	t.Parallel()
	app := New(nil)
	app.Router.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	rr := httptest.NewRecorder()

	assert.NotPanics(t, func() {
		app.Router.ServeHTTP(rr, req)
	})
}

func TestRouterRecoversFromPanicsReturns500(t *testing.T) {
	t.Parallel()
	app := New(nil)
	app.Router.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	rr := httptest.NewRecorder()

	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
