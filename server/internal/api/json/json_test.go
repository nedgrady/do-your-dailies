package apijson

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"do-your-dailies/server/internal/contracts"
	apperrors "do-your-dailies/server/internal/errors"

	"github.com/stretchr/testify/assert"
)

func TestDecodeBodyReturns422OnUnknownField(t *testing.T) {
	t.Parallel()
	body := strings.NewReader(`{"name":"dishes","unknown":1}`)
	req := httptest.NewRequest(http.MethodPost, "/api/chores/", body)
	var dst contracts.CreateChoreRequest

	err := DecodeBody(req, &dst)
	rr := httptest.NewRecorder()
	apperrors.Write(rr, err)

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}
