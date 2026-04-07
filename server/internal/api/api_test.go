package api

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"do-your-dailies/server/internal/models"
	"do-your-dailies/server/internal/store"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// mockChoreStore is a test double for store.ChoreStore.
type mockChoreStore struct {
	listFn   func() ([]models.Chore, error)
	createFn func(models.CreateChoreRequest) (models.Chore, error)
	getFn    func(uint) (models.Chore, error)
	updateFn func(uint, models.UpdateChoreRequest) (models.Chore, error)
	deleteFn func(uint) error
}

func (m *mockChoreStore) List() ([]models.Chore, error) { return m.listFn() }
func (m *mockChoreStore) Create(req models.CreateChoreRequest) (models.Chore, error) {
	return m.createFn(req)
}
func (m *mockChoreStore) Get(id uint) (models.Chore, error) { return m.getFn(id) }
func (m *mockChoreStore) Update(id uint, req models.UpdateChoreRequest) (models.Chore, error) {
	return m.updateFn(id, req)
}
func (m *mockChoreStore) Delete(id uint) error { return m.deleteFn(id) }

var _ store.ChoreStore = (*mockChoreStore)(nil)

func newAppWithStore(mock store.ChoreStore) *Application {
	app := New(nil)
	app.ChoreStore = mock
	return app
}

// recordingResponseWriter records whether WriteHeader was called before Write.
type recordingResponseWriter struct {
	header                 http.Header
	statusCode             int
	headerWritten          bool
	writeBeforeWriteHeader bool
	body                   bytes.Buffer
}

func newRecordingResponseWriter() *recordingResponseWriter {
	return &recordingResponseWriter{header: make(http.Header)}
}

func (w *recordingResponseWriter) Header() http.Header {
	return w.header
}

func (w *recordingResponseWriter) WriteHeader(statusCode int) {
	w.headerWritten = true
	w.statusCode = statusCode
}

func (w *recordingResponseWriter) Write(payload []byte) (int, error) {
	if !w.headerWritten {
		w.writeBeforeWriteHeader = true
	}

	return w.body.Write(payload)
}

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

	body := strings.NewReader(`{"name":"dishes","cadence_in_days":1}`)
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

	body := strings.NewReader(`{"name":"dishes","cadence_in_days":1}`)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/api/chores/", body))

	assert.Contains(t, rr.Body.String(), "dishes")
}

func TestCreateChoreDecodesCadenceInDaysFromSnakeCaseJSON(t *testing.T) {
	app := newAppWithStore(&mockChoreStore{
		createFn: func(req models.CreateChoreRequest) (models.Chore, error) {
			return models.Chore{Name: "ok", CadenceInDays: req.CadenceInDays}, nil
		},
	})

	body := strings.NewReader(`{"name":"dishes","cadence_in_days":7}`)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/api/chores/", body))

	assert.Contains(t, rr.Body.String(), `"cadence_in_days":7`)
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

	body := strings.NewReader(`{"cadence_in_days":0}`)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodPut, "/api/chores/1", body))

	assert.Equal(t, http.StatusOK, rr.Code)
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

	body := strings.NewReader(`{"name":"dishes","cadence_in_days":1}`)
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

	body := strings.NewReader(`{"name":"dishes","cadence_in_days":1}`)
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
