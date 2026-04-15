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

func TestWriteReturnsNotFoundMessageForNotFound(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()

	Write(rr, MapStoreError(gorm.ErrRecordNotFound))

	assert.Equal(t, "not found\n", rr.Body.String())
}

func TestWriteReturns500ForUnexpectedError(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()

	Write(rr, MapStoreError(errors.New("boom")))

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestWriteReturnsInternalMessageForUnexpectedError(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()

	Write(rr, MapStoreError(errors.New("boom")))

	assert.Equal(t, "internal server error\n", rr.Body.String())
}

func TestWriteReturns500ForNilError(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()

	Write(rr, nil)

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

func TestMapStoreErrorUnexpectedErrorUnwrapsOriginalCause(t *testing.T) {
	t.Parallel()
	original := errors.New("boom")

	err := MapStoreError(original)

	assert.ErrorIs(t, err, original)
}

func TestMapStoreErrorRecordNotFoundMatchesNotFoundCategory(t *testing.T) {
	t.Parallel()

	err := MapStoreError(gorm.ErrRecordNotFound)

	assert.True(t, errors.Is(err, errNotFound))
}

func TestMapStoreErrorIsNotFoundCategory(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()

	Write(rr, MapStoreError(gorm.ErrRecordNotFound))

	assert.Equal(t, "not found\n", rr.Body.String())
}

func TestMapStoreErrorUnexpectedErrorIsNotNotFound(t *testing.T) {
	t.Parallel()

	err := MapStoreError(errors.New("boom"))

	assert.False(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestBadRequestWithoutCauseUsesCategoryMessage(t *testing.T) {
	t.Parallel()

	err := BadRequest(nil)

	assert.Equal(t, "bad request", err.Error())
}

func TestBadRequestUnwrapsCause(t *testing.T) {
	t.Parallel()
	original := errors.New("bad json")

	err := BadRequest(original)

	assert.ErrorIs(t, err, original)
}

func TestInternalWithoutCauseUsesCategoryMessage(t *testing.T) {
	t.Parallel()

	err := Internal(nil)

	assert.Equal(t, "internal server error", err.Error())
}

func TestInternalUnwrapsCause(t *testing.T) {
	t.Parallel()
	original := errors.New("boom")

	err := Internal(original)

	assert.ErrorIs(t, err, original)
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

func TestBadRequestErrorIsDoesNotMatchNilTarget(t *testing.T) {
	t.Parallel()

	err := BadRequest(nil)

	assert.False(t, errors.Is(err, nil))
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

func TestInternalErrorIsMatchesInternalCategory(t *testing.T) {
	t.Parallel()

	err := Internal(nil)

	assert.True(t, errors.Is(err, errInternal))
}

func TestInternalErrorIsDoesNotMatchBadRequestCategory(t *testing.T) {
	t.Parallel()

	err := Internal(nil)

	assert.False(t, errors.Is(err, errBadRequest))
}
