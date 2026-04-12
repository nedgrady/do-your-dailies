package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"do-your-dailies/server/internal/domain/chorecompletions"
	"do-your-dailies/server/internal/domain/chores"
	"do-your-dailies/server/internal/testhelpers"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNewWithNilDBLeavesStoresNil(t *testing.T) {
	t.Parallel()

	app := New(nil)

	assert.Nil(t, app.ChoreStore)
}

func TestNewWithDBInitializesChoreStore(t *testing.T) {
	t.Parallel()
	database := testhelpers.NewTransactionDB(t, migrateForAPITests)

	app := New(database)

	assert.NotNil(t, app.ChoreStore)
}

func TestNewWithDBInitializesChoreCompletionStore(t *testing.T) {
	t.Parallel()
	database := testhelpers.NewTransactionDB(t, migrateForAPITests)

	app := New(database)

	assert.NotNil(t, app.ChoreCompletionStore)
}

func TestRouterRegistersChoresRoute(t *testing.T) {
	t.Parallel()
	app := New(nil)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/chores/", nil)

	app.Router.ServeHTTP(rr, req)

	assert.NotEqual(t, http.StatusNotFound, rr.Code)
}

func TestRouterRegistersChoreCompletionsRoute(t *testing.T) {
	t.Parallel()
	app := New(nil)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/chore-completions/", strings.NewReader(`not-json`))

	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestRouterRegistersOpenAPISpecRoute(t *testing.T) {
	t.Parallel()
	app := New(nil)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)

	app.Router.ServeHTTP(rr, req)

	assert.NotEqual(t, http.StatusNotFound, rr.Code)
}

func TestRouterRegistersSwaggerRedirectRoute(t *testing.T) {
	t.Parallel()
	app := New(nil)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/swagger", nil)

	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMovedPermanently, rr.Code)
}

func TestRouterRecoversWhenChoreStoreIsNil(t *testing.T) {
	t.Parallel()
	app := New(nil)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/chores/", nil)

	app.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func migrateForAPITests(database *gorm.DB) error {
	return database.AutoMigrate(&chores.Chore{}, &chorecompletions.ChoreCompletion{})
}
