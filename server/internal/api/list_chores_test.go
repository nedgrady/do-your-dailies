package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"do-your-dailies/server/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestListChoresReturns200(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{
		listFn: func() ([]models.Chore, error) { return []models.Chore{}, nil },
	})

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chores/", nil))

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestListChoresReturnsJSONContentType(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{
		listFn: func() ([]models.Chore, error) { return []models.Chore{}, nil },
	})

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chores/", nil))

	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}

func TestListChoresReturnsChores(t *testing.T) {
	chore := models.Chore{Name: "dishes", CadenceInDays: 1}
	app := newAppWithStore(&mockChoreStore{
		listFn: func() ([]models.Chore, error) { return []models.Chore{chore}, nil },
	})

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chores/", nil))

	assert.Contains(t, rr.Body.String(), "dishes")
}

func TestListChoresReturns500OnStoreError(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{
		listFn: func() ([]models.Chore, error) { return nil, errors.New("db error") },
	})

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chores/", nil))

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
