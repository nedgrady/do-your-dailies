package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"do-your-dailies/server/internal/models"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetChoreReturns200(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{
		getFn: func(id uint) (models.Chore, error) {
			return models.Chore{Name: "dishes"}, nil
		},
	})

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chores/1", nil))

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGetChoreReturnsChore(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{
		getFn: func(id uint) (models.Chore, error) {
			return models.Chore{Name: "dishes"}, nil
		},
	})

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chores/1", nil))

	assert.Contains(t, rr.Body.String(), "dishes")
}

func TestGetChoreReturnsLowerCamelCaseID(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{
		getFn: func(id uint) (models.Chore, error) {
			return models.Chore{Model: gorm.Model{ID: 42}, Name: "dishes"}, nil
		},
	})

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chores/42", nil))

	assert.Contains(t, rr.Body.String(), `"id":42`)
}

func TestGetChoreReturnsLowerCamelCaseCreatedAt(t *testing.T) {
	createdAt := time.Date(2026, time.April, 7, 12, 0, 0, 0, time.UTC)
	app := newAppWithStore(&mockChoreStore{
		getFn: func(id uint) (models.Chore, error) {
			return models.Chore{Model: gorm.Model{CreatedAt: createdAt}, Name: "dishes"}, nil
		},
	})

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chores/42", nil))

	assert.Contains(t, rr.Body.String(), `"createdAt":"2026-04-07T12:00:00Z"`)
}

func TestGetChoreReturnsLowerCamelCaseUpdatedAt(t *testing.T) {
	updatedAt := time.Date(2026, time.April, 7, 13, 0, 0, 0, time.UTC)
	app := newAppWithStore(&mockChoreStore{
		getFn: func(id uint) (models.Chore, error) {
			return models.Chore{Model: gorm.Model{UpdatedAt: updatedAt}, Name: "dishes"}, nil
		},
	})

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chores/42", nil))

	assert.Contains(t, rr.Body.String(), `"updatedAt":"2026-04-07T13:00:00Z"`)
}

func TestGetChoreReturns404WhenNotFound(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{
		getFn: func(id uint) (models.Chore, error) {
			return models.Chore{}, gorm.ErrRecordNotFound
		},
	})

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chores/99", nil))

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetChoreReturns400OnNonNumericID(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{})

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chores/abc", nil))

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetChoreReturns500OnUnexpectedError(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{
		getFn: func(id uint) (models.Chore, error) {
			return models.Chore{}, errors.New("db error")
		},
	})

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chores/1", nil))

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGetChoreAcceptsLargeUintID(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{
		getFn: func(id uint) (models.Chore, error) {
			return models.Chore{Name: "large-id"}, nil
		},
	})

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chores/9223372036854775808", nil))

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGetChoreRejectsBase11OnlyID(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{})

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chores/a", nil))

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
