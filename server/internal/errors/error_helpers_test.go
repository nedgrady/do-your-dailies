package errors

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
)

func TestWriteReturns404ForNotFound(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()

	Write(rr, MapStoreError(gorm.ErrRecordNotFound))

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestWriteReturns500ForUnexpectedError(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()

	Write(rr, MapStoreError(errors.New("boom")))

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestWriteReturns422ForBadRequestCategory(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()

	Write(rr, BadRequest(errors.New("bad json")))

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestMapStoreErrorUnwrapsRecordNotFoundCause(t *testing.T) {
	t.Parallel()

	err := MapStoreError(gorm.ErrRecordNotFound)

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestMapStoreErrorIsNotFoundCategory(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()

	Write(rr, MapStoreError(gorm.ErrRecordNotFound))

	assert.Equal(t, "not found\n", rr.Body.String())
}

func TestBadRequestWithoutCauseUsesCategoryMessage(t *testing.T) {
	t.Parallel()

	err := BadRequest(nil)

	assert.Equal(t, "bad request", err.Error())
}

func TestInternalWithoutCauseUsesCategoryMessage(t *testing.T) {
	t.Parallel()

	err := Internal(nil)

	assert.Equal(t, "internal server error", err.Error())
}

func TestBadRequestMatchesCategoryWithErrorIs(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()

	Write(rr, BadRequest(nil))

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestBadRequestWritesExpectedMessage(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()

	Write(rr, BadRequest(nil))

	assert.Equal(t, "bad request\n", rr.Body.String())
}

func TestBadRequestErrorIsMatchesBadRequestCategory(t *testing.T) {
	t.Parallel()

	err := BadRequest(nil)

	assert.True(t, errors.Is(err, errBadRequest))
}

func TestBadRequestErrorIsDoesNotMatchNotFoundCategory(t *testing.T) {
	t.Parallel()

	err := BadRequest(nil)

	assert.False(t, errors.Is(err, errNotFound))
}

func TestInternalMatchesCategoryWithErrorIs(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()

	Write(rr, Internal(nil))

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestInternalWritesExpectedMessage(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()

	Write(rr, Internal(nil))

	assert.Equal(t, "internal server error\n", rr.Body.String())
}
