package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteChoreReturns204(t *testing.T) {
	app, database := newPostgresTestApp(t)
	chore := seedChore(t, database, "dishes", 1)

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/chores/%d", chore.ID), nil))

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestDeleteChoreReturns404WhenNotFound(t *testing.T) {
	app, _ := newPostgresTestApp(t)

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodDelete, "/api/chores/99", nil))

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestDeleteChoreReturns400OnNonNumericID(t *testing.T) {
	app, _ := newPostgresTestApp(t)

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodDelete, "/api/chores/abc", nil))

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestDeleteChoreReturns500OnUnexpectedError(t *testing.T) {
	app, database := newPostgresTestApp(t)
	chore := seedChore(t, database, "dishes", 1)
	breakChoreTable(t, database)

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/chores/%d", chore.ID), nil))

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestDeleteChoreAcceptsLargeUintID(t *testing.T) {
	app, _ := newPostgresTestApp(t)

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodDelete, "/api/chores/9223372036854775808", nil))

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestDeleteChoreRejectsBase11OnlyID(t *testing.T) {
	app, _ := newPostgresTestApp(t)

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodDelete, "/api/chores/a", nil))

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
