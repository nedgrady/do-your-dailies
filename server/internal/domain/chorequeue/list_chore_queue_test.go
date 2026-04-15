package chorequeue

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"do-your-dailies/server/internal/domain/chores"

	"github.com/stretchr/testify/assert"
)

func TestListChoreQueueReturns200(t *testing.T) {
	t.Parallel()
	router, _ := newPostgresTestRouter(t)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chore-queue/", nil))

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestListChoreQueueReturnsJSONContentType(t *testing.T) {
	t.Parallel()
	router, _ := newPostgresTestRouter(t)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chore-queue/", nil))

	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}

func TestListChoreQueueReturnsOrderedChores(t *testing.T) {
	t.Parallel()
	router, database := newPostgresTestRouter(t)
	seedChoreWithCompletion(t, database, "bathroom", 30, -100, intPtr(-40))
	seedChoreWithCompletion(t, database, "dishes", 1, -100, intPtr(-1))
	seedChoreWithCompletion(t, database, "vacuum", 7, -100, intPtr(-8))

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chore-queue/", nil))

	assert.Equal(t, []string{"bathroom", "vacuum", "dishes"}, queueNamesFromBody(t, rr.Body.String()))
}

func TestListChoreQueueAppliesMaxChores(t *testing.T) {
	t.Parallel()
	router, database := newPostgresTestRouter(t)
	seedChoreWithCompletion(t, database, "bathroom", 30, -100, intPtr(-40))
	seedChoreWithCompletion(t, database, "dishes", 1, -100, intPtr(-1))
	seedChoreWithCompletion(t, database, "vacuum", 7, -100, intPtr(-8))
	seedChoreWithCompletion(t, database, "bins", 7, -100, intPtr(-20))

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chore-queue/?maxChores=2", nil))

	assert.Equal(t, []string{"bins", "bathroom"}, queueNamesFromBody(t, rr.Body.String()))
}

func TestListChoreQueueReturns400OnInvalidDayOffset(t *testing.T) {
	t.Parallel()
	router, _ := newPostgresTestRouter(t)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chore-queue/?dayOffset=tomorrow", nil))

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestListChoreQueueReturns400OnInvalidMaxChores(t *testing.T) {
	t.Parallel()
	router, _ := newPostgresTestRouter(t)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chore-queue/?maxChores=0", nil))

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestListChoreQueueReturns500OnStoreError(t *testing.T) {
	t.Parallel()
	router := newTestRouter(failingStore{}, fixedNow)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chore-queue/", nil))

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

type failingStore struct{}

func (failingStore) List(targetDay time.Time, maxChores int) ([]chores.Chore, error) {
	return nil, errors.New("boom")
}

func queueNamesFromBody(t *testing.T, body string) []string {
	t.Helper()

	var queue []struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal([]byte(body), &queue); err != nil {
		t.Fatalf("decode queue body: %v", err)
	}

	names := make([]string, 0, len(queue))
	for _, chore := range queue {
		names = append(names, chore.Name)
	}

	return names
}
