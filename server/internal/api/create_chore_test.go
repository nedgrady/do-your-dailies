package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"do-your-dailies/server/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestCreateChoreReturns201(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{
		createFn: func(req models.CreateChoreRequest) (models.Chore, error) {
			return models.Chore{Name: req.Name, CadenceInDays: req.CadenceInDays}, nil
		},
	})

	body := strings.NewReader(`{"name":"dishes","cadenceInDays":1}`)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/api/chores/", body))

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestCreateChoreReturnsCreatedChore(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{
		createFn: func(req models.CreateChoreRequest) (models.Chore, error) {
			return models.Chore{Name: req.Name, CadenceInDays: req.CadenceInDays}, nil
		},
	})

	body := strings.NewReader(`{"name":"dishes","cadenceInDays":1}`)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/api/chores/", body))

	assert.Contains(t, rr.Body.String(), "dishes")
}

func TestCreateChoreRejectsSnakeCaseCadenceJSON(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{
		createFn: func(req models.CreateChoreRequest) (models.Chore, error) {
			return models.Chore{Name: req.Name, CadenceInDays: req.CadenceInDays}, nil
		},
	})

	body := strings.NewReader(`{"name":"dishes","cadence_in_days":7}`)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/api/chores/", body))

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestCreateChoreDecodesCadenceInDaysFromLowerCamelJSON(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{
		createFn: func(req models.CreateChoreRequest) (models.Chore, error) {
			return models.Chore{Name: req.Name, CadenceInDays: req.CadenceInDays}, nil
		},
	})

	body := strings.NewReader(`{"name":"dishes","cadenceInDays":7}`)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/api/chores/", body))

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestCreateChoreReturnsLowerCamelCaseCadence(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{
		createFn: func(req models.CreateChoreRequest) (models.Chore, error) {
			return models.Chore{Name: req.Name, CadenceInDays: req.CadenceInDays}, nil
		},
	})

	body := strings.NewReader(`{"name":"dishes","cadenceInDays":7}`)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/api/chores/", body))

	assert.Contains(t, rr.Body.String(), `"cadenceInDays":7`)
}

func TestCreateChoreReturns422OnBadJSON(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{})

	body := strings.NewReader(`not-json`)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/api/chores/", body))

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestCreateChoreReturns500OnStoreError(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{
		createFn: func(req models.CreateChoreRequest) (models.Chore, error) {
			return models.Chore{}, errors.New("db error")
		},
	})

	body := strings.NewReader(`{"name":"dishes","cadenceInDays":1}`)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/api/chores/", body))

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestCreateChoreStoreErrorBody(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{
		createFn: func(req models.CreateChoreRequest) (models.Chore, error) {
			return models.Chore{}, errors.New("db error")
		},
	})

	body := strings.NewReader(`{"name":"dishes","cadenceInDays":1}`)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/api/chores/", body))

	assert.Equal(t, "internal server error\n", rr.Body.String())
}

func TestCreateChoreWithRealDBReturns201(t *testing.T) {
	app := New(newCreateChoreTestDB(t))

	body := strings.NewReader(`{"name":"dishes","cadenceInDays":1}`)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/api/chores/", body))

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestCreateChoreWithRealDBPersistsChore(t *testing.T) {
	app := New(newCreateChoreTestDB(t))

	createBody := strings.NewReader(`{"name":"dishes","cadenceInDays":1}`)
	app.Router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/api/chores/", createBody))

	listResponse := httptest.NewRecorder()
	app.Router.ServeHTTP(listResponse, httptest.NewRequest(http.MethodGet, "/api/chores/", nil))

	assert.Contains(t, listResponse.Body.String(), "dishes")
}

func TestCreateChoreWithTransactionDBReturns201(t *testing.T) {
	app := New(newCreateChoreTestTransactionDB(t))

	body := strings.NewReader(`{"name":"dishes","cadenceInDays":1}`)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/api/chores/", body))

	assert.Equal(t, http.StatusCreated, rr.Code)
}
