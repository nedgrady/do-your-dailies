package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDeleteChoreReturns204(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{
		deleteFn: func(id uint) error { return nil },
	})

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodDelete, "/api/chores/1", nil))

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestDeleteChoreReturns404WhenNotFound(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{
		deleteFn: func(id uint) error { return gorm.ErrRecordNotFound },
	})

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodDelete, "/api/chores/99", nil))

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestDeleteChoreReturns400OnNonNumericID(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{})

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodDelete, "/api/chores/abc", nil))

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestDeleteChoreReturns500OnUnexpectedError(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{
		deleteFn: func(id uint) error { return errors.New("db error") },
	})

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodDelete, "/api/chores/1", nil))

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestDeleteChoreAcceptsLargeUintID(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{
		deleteFn: func(id uint) error { return nil },
	})

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodDelete, "/api/chores/9223372036854775808", nil))

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestDeleteChoreRejectsBase11OnlyID(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{})

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodDelete, "/api/chores/a", nil))

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
