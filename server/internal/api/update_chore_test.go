package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"do-your-dailies/server/internal/models"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUpdateChoreReturns200(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{
		updateFn: func(id uint, req models.UpdateChoreRequest) (models.Chore, error) {
			return models.Chore{Name: *req.Name}, nil
		},
	})

	body := strings.NewReader(`{"name":"vacuuming"}`)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodPut, "/api/chores/1", body))

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestUpdateChoreReturnsUpdatedChore(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{
		updateFn: func(id uint, req models.UpdateChoreRequest) (models.Chore, error) {
			return models.Chore{Name: *req.Name}, nil
		},
	})

	body := strings.NewReader(`{"name":"vacuuming"}`)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodPut, "/api/chores/1", body))

	assert.Contains(t, rr.Body.String(), "vacuuming")
}

func TestUpdateChoreOmitsCadenceAsNilPointer(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{
		updateFn: func(id uint, req models.UpdateChoreRequest) (models.Chore, error) {
			if req.CadenceInDays != nil {
				return models.Chore{}, errors.New("cadence should be nil when omitted")
			}
			return models.Chore{Name: "ok"}, nil
		},
	})

	body := strings.NewReader(`{"name":"vacuuming"}`)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodPut, "/api/chores/1", body))

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestUpdateChoreAllowsExplicitZeroCadence(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{
		updateFn: func(id uint, req models.UpdateChoreRequest) (models.Chore, error) {
			if req.CadenceInDays == nil {
				return models.Chore{}, errors.New("cadence should be set")
			}
			if *req.CadenceInDays != 0 {
				return models.Chore{}, errors.New("cadence should decode zero")
			}
			return models.Chore{CadenceInDays: *req.CadenceInDays}, nil
		},
	})

	body := strings.NewReader(`{"cadenceInDays":0}`)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodPut, "/api/chores/1", body))

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestUpdateChoreDecodesLowerCamelCaseCadence(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{
		updateFn: func(id uint, req models.UpdateChoreRequest) (models.Chore, error) {
			if req.CadenceInDays == nil {
				return models.Chore{}, errors.New("cadence should decode")
			}
			return models.Chore{CadenceInDays: *req.CadenceInDays}, nil
		},
	})

	body := strings.NewReader(`{"cadenceInDays":3}`)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodPut, "/api/chores/1", body))

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestUpdateChoreRejectsSnakeCaseCadenceJSON(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{
		updateFn: func(id uint, req models.UpdateChoreRequest) (models.Chore, error) {
			return models.Chore{Name: "unexpected"}, nil
		},
	})

	body := strings.NewReader(`{"cadence_in_days":3}`)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodPut, "/api/chores/1", body))

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestUpdateChoreReturnsLowerCamelCaseCadence(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{
		updateFn: func(id uint, req models.UpdateChoreRequest) (models.Chore, error) {
			return models.Chore{CadenceInDays: 3}, nil
		},
	})

	body := strings.NewReader(`{"cadenceInDays":3}`)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodPut, "/api/chores/1", body))

	assert.Contains(t, rr.Body.String(), `"cadenceInDays":3`)
}

func TestUpdateChoreReturns404WhenNotFound(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{
		updateFn: func(id uint, req models.UpdateChoreRequest) (models.Chore, error) {
			return models.Chore{}, gorm.ErrRecordNotFound
		},
	})

	body := strings.NewReader(`{"name":"vacuuming"}`)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodPut, "/api/chores/99", body))

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestUpdateChoreReturns400OnNonNumericID(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{})

	body := strings.NewReader(`{"name":"vacuuming"}`)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodPut, "/api/chores/abc", body))

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUpdateChoreReturns422OnBadJSON(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{})

	body := strings.NewReader(`not-json`)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodPut, "/api/chores/1", body))

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestUpdateChoreReturns500OnUnexpectedError(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{
		updateFn: func(id uint, req models.UpdateChoreRequest) (models.Chore, error) {
			return models.Chore{}, errors.New("db error")
		},
	})

	body := strings.NewReader(`{"name":"vacuuming"}`)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodPut, "/api/chores/1", body))

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestUpdateChoreAcceptsLargeUintID(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{
		updateFn: func(id uint, req models.UpdateChoreRequest) (models.Chore, error) {
			return models.Chore{Name: "large-id"}, nil
		},
	})

	body := strings.NewReader(`{"name":"vacuuming"}`)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodPut, "/api/chores/9223372036854775808", body))

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestUpdateChoreRejectsBase11OnlyID(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{})

	body := strings.NewReader(`{"name":"vacuuming"}`)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodPut, "/api/chores/a", body))

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
