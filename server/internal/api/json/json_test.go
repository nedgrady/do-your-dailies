package apijson

import (
	"encoding/json"
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

func TestDecodeBodyDecodesValidPayload(t *testing.T) {
	t.Parallel()
	body := strings.NewReader(`{"name":"dishes","cadenceInDays":3}`)
	req := httptest.NewRequest(http.MethodPost, "/api/chores/", body)
	var dst contracts.CreateChoreRequest

	err := DecodeBody(req, &dst)

	assert.NoError(t, err)
}

func TestDecodeBodyFillsDestinationFields(t *testing.T) {
	t.Parallel()
	body := strings.NewReader(`{"name":"dishes","cadenceInDays":3}`)
	req := httptest.NewRequest(http.MethodPost, "/api/chores/", body)
	var dst contracts.CreateChoreRequest

	_ = DecodeBody(req, &dst)

	assert.Equal(t, "dishes", dst.Name)
}

func TestWriteSetsJSONContentType(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()

	Write(rr, http.StatusCreated, map[string]string{"name": "dishes"})

	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}

func TestWriteUsesProvidedStatus(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()

	Write(rr, http.StatusCreated, map[string]string{"name": "dishes"})

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestWriteEncodesPayload(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()

	Write(rr, http.StatusCreated, map[string]string{"name": "dishes"})

	var payload map[string]string
	_ = json.Unmarshal(rr.Body.Bytes(), &payload)
	assert.Equal(t, "dishes", payload["name"])
}
