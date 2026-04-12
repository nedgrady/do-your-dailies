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
