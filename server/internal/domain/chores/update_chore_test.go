package chores

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateChoreReturns200(t *testing.T) {
	t.Parallel()
	router, database := newPostgresTestRouter(t)
	chore := seedChore(t, database, "dishes", 1)

	body := strings.NewReader(`{"name":"vacuuming"}`)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/chores/%d", chore.ID), body))

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestUpdateChoreReturnsUpdatedChore(t *testing.T) {
	t.Parallel()
	router, database := newPostgresTestRouter(t)
	chore := seedChore(t, database, "dishes", 1)

	body := strings.NewReader(`{"name":"vacuuming"}`)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/chores/%d", chore.ID), body))

	assert.Contains(t, rr.Body.String(), "vacuuming")
}

func TestUpdateChoreOmitsCadenceAsNilPointer(t *testing.T) {
	t.Parallel()
	router, database := newPostgresTestRouter(t)
	chore := seedChore(t, database, "dishes", 7)

	body := strings.NewReader(`{"name":"vacuuming"}`)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/chores/%d", chore.ID), body))

	assert.Contains(t, rr.Body.String(), `"cadenceInDays":7`)
}

func TestUpdateChoreRejectsExplicitZeroCadence(t *testing.T) {
	t.Parallel()
	router, database := newPostgresTestRouter(t)
	chore := seedChore(t, database, "dishes", 7)

	body := strings.NewReader(`{"cadenceInDays":0}`)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/chores/%d", chore.ID), body))

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestUpdateChoreDecodesLowerCamelCaseCadence(t *testing.T) {
	t.Parallel()
	router, database := newPostgresTestRouter(t)
	chore := seedChore(t, database, "dishes", 1)

	body := strings.NewReader(`{"cadenceInDays":3}`)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/chores/%d", chore.ID), body))

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestUpdateChoreRejectsSnakeCaseCadenceJSON(t *testing.T) {
	t.Parallel()
	router, database := newPostgresTestRouter(t)
	chore := seedChore(t, database, "dishes", 1)

	body := strings.NewReader(`{"cadence_in_days":3}`)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/chores/%d", chore.ID), body))

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestUpdateChoreReturnsLowerCamelCaseCadence(t *testing.T) {
	t.Parallel()
	router, database := newPostgresTestRouter(t)
	chore := seedChore(t, database, "dishes", 1)

	body := strings.NewReader(`{"cadenceInDays":3}`)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/chores/%d", chore.ID), body))

	assert.Contains(t, rr.Body.String(), `"cadenceInDays":3`)
}

func TestUpdateChoreReturns404WhenNotFound(t *testing.T) {
	t.Parallel()
	router, _ := newPostgresTestRouter(t)

	body := strings.NewReader(`{"name":"vacuuming"}`)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodPut, "/api/chores/99", body))

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestUpdateChoreReturns400OnNonNumericID(t *testing.T) {
	t.Parallel()
	router, _ := newPostgresTestRouter(t)

	body := strings.NewReader(`{"name":"vacuuming"}`)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodPut, "/api/chores/abc", body))

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUpdateChoreReturns422OnBadJSON(t *testing.T) {
	t.Parallel()
	router, _ := newPostgresTestRouter(t)

	body := strings.NewReader(`not-json`)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodPut, "/api/chores/1", body))

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestUpdateChoreReturns500OnUnexpectedError(t *testing.T) {
	t.Parallel()
	router, database := newPostgresTestRouter(t)
	chore := seedChore(t, database, "dishes", 1)
	breakChoreTable(t, database)

	body := strings.NewReader(`{"name":"vacuuming"}`)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/chores/%d", chore.ID), body))

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestUpdateChoreAcceptsLargeUintID(t *testing.T) {
	t.Parallel()
	router, _ := newPostgresTestRouter(t)

	body := strings.NewReader(`{"name":"vacuuming"}`)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodPut, "/api/chores/9223372036854775808", body))

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestUpdateChoreRejectsBase11OnlyID(t *testing.T) {
	t.Parallel()
	router, _ := newPostgresTestRouter(t)

	body := strings.NewReader(`{"name":"vacuuming"}`)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodPut, "/api/chores/a", body))

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
