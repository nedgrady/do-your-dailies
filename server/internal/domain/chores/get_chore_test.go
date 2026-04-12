package chores

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetChoreReturns200(t *testing.T) {
	t.Parallel()
	router, database := newPostgresTestRouter(t)
	chore := seedChore(t, database, "dishes", 1)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/chores/%d", chore.ID), nil))

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGetChoreReturnsChore(t *testing.T) {
	t.Parallel()
	router, database := newPostgresTestRouter(t)
	chore := seedChore(t, database, "dishes", 1)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/chores/%d", chore.ID), nil))

	assert.Contains(t, rr.Body.String(), "dishes")
}

func TestGetChoreReturnsLowerCamelCaseID(t *testing.T) {
	t.Parallel()
	router, database := newPostgresTestRouter(t)
	chore := seedChore(t, database, "dishes", 1)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/chores/%d", chore.ID), nil))

	assert.Contains(t, rr.Body.String(), fmt.Sprintf(`"id":%d`, chore.ID))
}

func TestGetChoreReturnsLowerCamelCaseCreatedAt(t *testing.T) {
	t.Parallel()
	createdAt := time.Date(2026, time.April, 7, 12, 0, 0, 0, time.UTC)
	router, database := newPostgresTestRouter(t)
	chore := seedChore(t, database, "dishes", 1)
	if err := database.Model(&chore).Update("created_at", createdAt).Error; err != nil {
		t.Fatalf("set created_at: %v", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/chores/%d", chore.ID), nil))

	assert.Contains(t, rr.Body.String(), `"createdAt":"`)
}

func TestGetChoreReturnsLowerCamelCaseUpdatedAt(t *testing.T) {
	t.Parallel()
	updatedAt := time.Date(2026, time.April, 7, 13, 0, 0, 0, time.UTC)
	router, database := newPostgresTestRouter(t)
	chore := seedChore(t, database, "dishes", 1)
	if err := database.Model(&chore).Update("updated_at", updatedAt).Error; err != nil {
		t.Fatalf("set updated_at: %v", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/chores/%d", chore.ID), nil))

	assert.Contains(t, rr.Body.String(), `"updatedAt":"`)
}

func TestGetChoreReturns404WhenNotFound(t *testing.T) {
	t.Parallel()
	router, _ := newPostgresTestRouter(t)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chores/99", nil))

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetChoreReturns400OnNonNumericID(t *testing.T) {
	t.Parallel()
	router, _ := newPostgresTestRouter(t)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chores/abc", nil))

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetChoreReturns500OnUnexpectedError(t *testing.T) {
	t.Parallel()
	router, database := newPostgresTestRouter(t)
	breakChoreTable(t, database)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chores/1", nil))

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGetChoreAcceptsLargeUintID(t *testing.T) {
	t.Parallel()
	router, _ := newPostgresTestRouter(t)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chores/9223372036854775808", nil))

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGetChoreRejectsBase11OnlyID(t *testing.T) {
	t.Parallel()
	router, _ := newPostgresTestRouter(t)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/chores/a", nil))

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
