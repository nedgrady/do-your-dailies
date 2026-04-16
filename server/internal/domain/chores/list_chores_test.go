package chores

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListChoresReturns200(t *testing.T) {
	t.Parallel()
	router, _ := newPostgresTestRouter(t)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chores/", nil))

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestListChoresReturnsJSONContentType(t *testing.T) {
	t.Parallel()
	router, _ := newPostgresTestRouter(t)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chores/", nil))

	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}

func TestListChoresReturnsChores(t *testing.T) {
	t.Parallel()
	router, database := newPostgresTestRouter(t)
	seedChore(t, database, "dishes", 1)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chores/", nil))

	assert.Contains(t, rr.Body.String(), "dishes")
}

func TestListChoresReturns500OnStoreError(t *testing.T) {
	t.Parallel()
	router := newTestRouter(failingStore{listErr: errors.New("boom")})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chores/", nil))

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
