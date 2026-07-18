package apijson

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

func TestDecodeBodyNormalizesDateTimeToUTC(t *testing.T) {
	t.Parallel()
	body := strings.NewReader(`{"id":1,"name":"dishes","cadenceInDays":3,"createdAt":"2026-04-07T12:00:00+02:00","updatedAt":"2026-04-07T12:00:00+02:00"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/chores/", body)
	var dst contracts.Chore

	if err := DecodeBody(req, &dst); err != nil {
		t.Fatalf("decode body: %v", err)
	}

	assert.Equal(t, "2026-04-07T10:00:00Z", dst.CreatedAt.Format(time.RFC3339))
}

func TestWriteEmitsDateTimeAsUTCZulu(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()
	payload := contracts.Chore{
		Id:            1,
		Name:          "dishes",
		CadenceInDays: 3,
		CreatedAt:     time.Date(2026, time.April, 7, 12, 0, 0, 0, time.FixedZone("UTC+2", 2*60*60)),
		UpdatedAt:     time.Date(2026, time.April, 7, 12, 0, 0, 0, time.FixedZone("UTC+2", 2*60*60)),
	}

	Write(rr, http.StatusOK, payload)

	assert.Contains(t, rr.Body.String(), `"createdAt":"2026-04-07T10:00:00Z"`)
}
