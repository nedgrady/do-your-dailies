package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListChoresReturns200(t *testing.T) {
	t.Parallel()
	app, _ := newPostgresTestApp(t)

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chores/", nil))

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestListChoresReturnsJSONContentType(t *testing.T) {
	t.Parallel()
	app, _ := newPostgresTestApp(t)

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chores/", nil))

	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}

func TestListChoresReturnsChores(t *testing.T) {
	t.Parallel()
	app, database := newPostgresTestApp(t)
	seedChore(t, database, "dishes", 1)

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chores/", nil))

	assert.Contains(t, rr.Body.String(), "dishes")
}

func TestListChoresReturns500OnStoreError(t *testing.T) {
	t.Parallel()
	app, database := newPostgresTestApp(t)
	breakChoreTable(t, database)

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chores/", nil))

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
