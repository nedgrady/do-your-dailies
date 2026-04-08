package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"do-your-dailies/server/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestCreateChoreCompletionReturns201(t *testing.T) {
	t.Parallel()
	app, database := newPostgresTestApp(t)
	chore := seedChore(t, database, "dishes", 1)

	body := strings.NewReader(`{"choreId":` + itoa(chore.ID) + `}`)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/api/chore-completions/", body))

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestCreateChoreCompletionReturnsCreatedCompletion(t *testing.T) {
	t.Parallel()
	app, database := newPostgresTestApp(t)
	chore := seedChore(t, database, "dishes", 1)

	body := strings.NewReader(`{"choreId":` + itoa(chore.ID) + `}`)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/api/chore-completions/", body))

	assert.Contains(t, rr.Body.String(), `"choreId":`+itoa(chore.ID))
}

func TestCreateChoreCompletionRejectsSnakeCaseChoreIDJSON(t *testing.T) {
	t.Parallel()
	app, database := newPostgresTestApp(t)
	chore := seedChore(t, database, "dishes", 1)

	body := strings.NewReader(`{"chore_id":` + itoa(chore.ID) + `}`)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/api/chore-completions/", body))

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestCreateChoreCompletionReturns422OnBadJSON(t *testing.T) {
	t.Parallel()
	app, _ := newPostgresTestApp(t)

	body := strings.NewReader(`not-json`)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/api/chore-completions/", body))

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestCreateChoreCompletionReturns500OnStoreError(t *testing.T) {
	t.Parallel()
	app, database := newPostgresTestApp(t)
	chore := seedChore(t, database, "dishes", 1)
	app.ChoreCompletionStore = failingChoreCompletionStore{}

	body := strings.NewReader(`{"choreId":` + itoa(chore.ID) + `}`)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/api/chore-completions/", body))

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestCreateChoreCompletionStoreErrorBody(t *testing.T) {
	t.Parallel()
	app, database := newPostgresTestApp(t)
	chore := seedChore(t, database, "dishes", 1)
	app.ChoreCompletionStore = failingChoreCompletionStore{}

	body := strings.NewReader(`{"choreId":` + itoa(chore.ID) + `}`)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/api/chore-completions/", body))

	assert.Equal(t, "internal server error\n", rr.Body.String())
}

func itoa(input uint) string {
	return strconv.FormatUint(uint64(input), 10)
}

type failingChoreCompletionStore struct{}

func (failingChoreCompletionStore) Create(models.CreateChoreCompletionRequest) (models.ChoreCompletion, error) {
	return models.ChoreCompletion{}, errors.New("boom")
}
