package api

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"do-your-dailies/server/internal/models"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestHealthCheck(t *testing.T) {
	app := New(nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestHealthCheckContentType(t *testing.T) {
	app := New(nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
}

func TestHealthCheckBody(t *testing.T) {
	app := New(nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, "healthy", rr.Body.String())
}

func TestOpenAPISpecRouteReturns200(t *testing.T) {
	app := New(nil)

	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	rr := httptest.NewRecorder()

	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestOpenAPISpecRouteReturnsSpecBody(t *testing.T) {
	app := New(nil)

	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	rr := httptest.NewRecorder()

	app.Router.ServeHTTP(rr, req)

	assert.Contains(t, rr.Body.String(), "openapi: 3.0.3")
}

func TestOpenAPISpecRouteReturnsYAMLContentType(t *testing.T) {
	app := New(nil)

	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	rr := httptest.NewRecorder()

	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, "application/yaml", rr.Header().Get("Content-Type"))
}

func TestSwaggerUIRouteReturns200(t *testing.T) {
	app := New(nil)

	req := httptest.NewRequest(http.MethodGet, "/swagger/", nil)
	rr := httptest.NewRecorder()

	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestSwaggerUIRouteReferencesSpec(t *testing.T) {
	app := New(nil)

	req := httptest.NewRequest(http.MethodGet, "/swagger/", nil)
	rr := httptest.NewRecorder()

	app.Router.ServeHTTP(rr, req)

	assert.Contains(t, rr.Body.String(), "/openapi.yaml")
}

func TestSwaggerRouteRedirects(t *testing.T) {
	app := New(nil)

	req := httptest.NewRequest(http.MethodGet, "/swagger", nil)
	rr := httptest.NewRecorder()

	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMovedPermanently, rr.Code)
}

func TestSwaggerRouteRedirectsToTrailingSlash(t *testing.T) {
	app := New(nil)

	req := httptest.NewRequest(http.MethodGet, "/swagger", nil)
	rr := httptest.NewRecorder()

	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, "/swagger/", rr.Header().Get("Location"))
}

func TestSwaggerUIRouteReturnsHTMLContentType(t *testing.T) {
	app := New(nil)

	req := httptest.NewRequest(http.MethodGet, "/swagger/", nil)
	rr := httptest.NewRecorder()

	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, "text/html; charset=utf-8", rr.Header().Get("Content-Type"))
}

func TestHealthCheckWritesExplicitStatusBeforeBody(t *testing.T) {
	app := New(nil)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	writer := newRecordingResponseWriter()

	app.healthCheck(writer, req)

	assert.False(t, writer.writeBeforeWriteHeader)
}

func TestHealthCheckWriterStatusCode(t *testing.T) {
	app := New(nil)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	writer := newRecordingResponseWriter()

	app.healthCheck(writer, req)

	assert.Equal(t, http.StatusOK, writer.statusCode)
}

func TestHealthCheckWriterContentType(t *testing.T) {
	app := New(nil)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	writer := newRecordingResponseWriter()

	app.healthCheck(writer, req)

	assert.Equal(t, "text/plain; charset=utf-8", writer.Header().Get("Content-Type"))
}

func TestHealthCheckWriterBody(t *testing.T) {
	app := New(nil)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	writer := newRecordingResponseWriter()

	app.healthCheck(writer, req)

	assert.Equal(t, "healthy", writer.body.String())
}

func TestRouterLogsHealthRequests(t *testing.T) {
	originalLogger := middleware.DefaultLogger
	var logOutput bytes.Buffer
	middleware.DefaultLogger = middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger:  log.New(&logOutput, "", 0),
		NoColor: true,
	})
	t.Cleanup(func() {
		middleware.DefaultLogger = originalLogger
	})

	app := New(nil)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	app.Router.ServeHTTP(rr, req)

	assert.NotEmpty(t, logOutput.String())
}

func TestRouterRecoversFromPanics(t *testing.T) {
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
	app := New(nil)
	app.Router.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	rr := httptest.NewRecorder()

	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// listChores

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

// createChore

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

// getChore

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

// updateChore

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

// deleteChore

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

func TestNewWithNonNilDBCreatesChoreStore(t *testing.T) {
	app := New(&gorm.DB{})

	assert.NotNil(t, app.ChoreStore)
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

func TestDeleteChoreReturns500OnUnexpectedError(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{
		deleteFn: func(id uint) error { return errors.New("db error") },
	})

	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodDelete, "/api/chores/1", nil))

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
