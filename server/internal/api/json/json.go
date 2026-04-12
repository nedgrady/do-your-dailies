package apijson

import (
	stdjson "encoding/json"
	"net/http"

	apperrors "do-your-dailies/server/internal/errors"
)

func DecodeBody(r *http.Request, dst any) error {
	decoder := stdjson.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return apperrors.BadRequest(err)
	}

	return nil
}

func Write(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = stdjson.NewEncoder(w).Encode(payload)
}
