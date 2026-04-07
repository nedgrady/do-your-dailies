package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeJSONBodyReturnsBadRequestErrorOnUnknownField(t *testing.T) {
	body := strings.NewReader(`{"name":"dishes","unknown":1}`)
	req := httptest.NewRequest(http.MethodPost, "/api/chores/", body)
	var dst CreateChoreRequest

	err := decodeJSONBody(req, &dst)

	assert.ErrorIs(t, err, errBadRequest)
}

func TestWriteAPIErrorReturns404ForNotFound(t *testing.T) {
	rr := httptest.NewRecorder()

	writeAPIError(rr, errNotFound)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestWriteAPIErrorReturns500ForUnexpectedError(t *testing.T) {
	rr := httptest.NewRecorder()

	writeAPIError(rr, errInternal)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
