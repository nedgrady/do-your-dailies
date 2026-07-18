package chorecompletions

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"do-your-dailies/server/internal/domain/choreinqueuecompletion"
	chorequeue "do-your-dailies/server/internal/domain/chorequeue"
	"do-your-dailies/server/internal/domain/chores"
	"do-your-dailies/server/internal/domain/models"
	"do-your-dailies/server/internal/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestCreateChoreCompletionReturns201(t *testing.T) {
	t.Parallel()
	router, database := newPostgresTestRouter(t)
	chore := seedChore(t, database, "dishes", 1)

	body := strings.NewReader(`{"choreId":` + itoa(chore.ID) + `}`)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/api/chore-completions/", body))

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestCreateChoreCompletionReturnsCreatedCompletion(t *testing.T) {
	t.Parallel()
	router, database := newPostgresTestRouter(t)
	chore := seedChore(t, database, "dishes", 1)

	body := strings.NewReader(`{"choreId":` + itoa(chore.ID) + `}`)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/api/chore-completions/", body))

	assert.Contains(t, rr.Body.String(), `"choreId":`+itoa(chore.ID))
}

func TestCreateChoreCompletionRejectsSnakeCaseChoreIDJSON(t *testing.T) {
	t.Parallel()
	router, database := newPostgresTestRouter(t)
	chore := seedChore(t, database, "dishes", 1)

	body := strings.NewReader(`{"chore_id":` + itoa(chore.ID) + `}`)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/api/chore-completions/", body))

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestCreateChoreCompletionReturns422OnBadJSON(t *testing.T) {
	t.Parallel()
	router, _ := newPostgresTestRouter(t)

	body := strings.NewReader(`not-json`)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/api/chore-completions/", body))

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestCreateChoreCompletionReturns500OnStoreError(t *testing.T) {
	t.Parallel()
	database := testhelpers.NewTransactionDB(t, migrate)
	chore := seedChore(t, database, "dishes", 1)
	router := newTestRouter(chores.NewGormStore(database), failingChoreCompletionStore{}, chorequeue.NewGormStore(database), choreinqueuecompletion.NewGormStore(database), database)

	body := strings.NewReader(`{"choreId":` + itoa(chore.ID) + `}`)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/api/chore-completions/", body))

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestCreateChoreCompletionStoreErrorBody(t *testing.T) {
	t.Parallel()
	database := testhelpers.NewTransactionDB(t, migrate)
	chore := seedChore(t, database, "dishes", 1)
	router := newTestRouter(chores.NewGormStore(database), failingChoreCompletionStore{}, chorequeue.NewGormStore(database), choreinqueuecompletion.NewGormStore(database), database)

	body := strings.NewReader(`{"choreId":` + itoa(chore.ID) + `}`)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/api/chore-completions/", body))

	assert.Equal(t, "internal server error\n", rr.Body.String())
}

func TestCreateChoreCompletionReturns404WhenChoreDoesNotExist(t *testing.T) {
	t.Parallel()
	router, _ := newPostgresTestRouter(t)

	body := strings.NewReader(`{"choreId":999999999}`)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/api/chore-completions/", body))

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

type failingChoreCompletionStore struct{}

func (failingChoreCompletionStore) Create(CreateRequest) (models.ChoreCompletion, error) {
	return models.ChoreCompletion{}, errors.New("boom")
}

func (failingChoreCompletionStore) ListByDay(time.Time) ([]models.ChoreCompletion, error) {
	return nil, errors.New("boom")
}
