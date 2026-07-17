package chorecompletions

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"do-your-dailies/server/internal/domain/choreinqueuecompletion"
	"do-your-dailies/server/internal/domain/chores"
	"do-your-dailies/server/internal/domain/models"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestListmodesReturns200(t *testing.T) {
	t.Parallel()
	router, _ := newPostgresTestRouter(t)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chore-completions/?date=2026-04-15", nil))

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestListmodesReturnsJSONContentType(t *testing.T) {
	t.Parallel()
	router, _ := newPostgresTestRouter(t)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chore-completions/?date=2026-04-15", nil))

	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}

func TestListmodesFiltersByDate(t *testing.T) {
	t.Parallel()
	router, database := newPostgresTestRouter(t)
	chore := seedChore(t, database, "dishes", 1)
	seedCompletion(t, database, chore.ID, time.Date(2026, time.April, 15, 10, 0, 0, 0, time.UTC))
	seedCompletion(t, database, chore.ID, time.Date(2026, time.April, 14, 10, 0, 0, 0, time.UTC))

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chore-completions/?date=2026-04-15", nil))

	assert.Equal(t, []uint64{uint64(chore.ID)}, choreIDsFromCompletionsBody(t, rr.Body.String()))
}

func TestListmodesDefaultsToTodayUTC(t *testing.T) {
	t.Parallel()
	router, database := newPostgresTestRouter(t)
	choreToday := seedChore(t, database, "dishes", 1)
	choreYesterday := seedChore(t, database, "bathroom", 7)
	now := time.Now().UTC()
	seedCompletion(t, database, choreToday.ID, time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, time.UTC))
	yesterday := now.AddDate(0, 0, -1)
	seedCompletion(t, database, choreYesterday.ID, time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 9, 0, 0, 0, time.UTC))

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chore-completions/", nil))

	assert.Equal(t, []uint64{uint64(choreToday.ID)}, choreIDsFromCompletionsBody(t, rr.Body.String()))
}

func TestListmodesReturns400OnInvalidDate(t *testing.T) {
	t.Parallel()
	router, _ := newPostgresTestRouter(t)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chore-completions/?date=15-04-2026", nil))

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestListmodesReturns500OnStoreError(t *testing.T) {
	t.Parallel()
	_, database := newPostgresTestRouter(t)
	router := newTestRouter(chores.NewGormStore(database), failingListStore{}, choreinqueuecompletion.NewGormStore(database), database)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chore-completions/?date=2026-04-15", nil))

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func choreIDsFromCompletionsBody(t *testing.T, body string) []uint64 {
	t.Helper()

	var completions []struct {
		ChoreID uint64 `json:"choreId"`
	}
	if err := json.Unmarshal([]byte(body), &completions); err != nil {
		t.Fatalf("decode completions body: %v", err)
	}

	ids := make([]uint64, 0, len(completions))
	for _, completion := range completions {
		ids = append(ids, completion.ChoreID)
	}

	return ids
}

func seedCompletion(t *testing.T, database *gorm.DB, choreID uint, completedAt time.Time) {
	t.Helper()

	completion := models.ChoreCompletion{ChoreID: choreID}
	if err := database.Create(&completion).Error; err != nil {
		t.Fatalf("seed chore completion: %v", err)
	}
	if err := database.Model(&completion).Updates(map[string]any{
		"created_at": completedAt,
		"updated_at": completedAt,
	}).Error; err != nil {
		t.Fatalf("set completion timestamps: %v", err)
	}
}

type failingListStore struct{}

func (failingListStore) Create(CreateRequest) (models.ChoreCompletion, error) {
	return models.ChoreCompletion{}, nil
}

func (failingListStore) ListByDay(time.Time) ([]models.ChoreCompletion, error) {
	return nil, errors.New("boom")
}
