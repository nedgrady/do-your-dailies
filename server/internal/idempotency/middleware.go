package idempotency

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	apijson "do-your-dailies/server/internal/api/json"
	"do-your-dailies/server/internal/auth"
	"do-your-dailies/server/internal/domain/models"

	"gorm.io/gorm"
)

const (
	HeaderName     = "Idempotency-Key"
	replayedHeader = "Idempotent-Replayed"
	maxBodyBytes   = 1 << 20
	maxKeyLength   = 64
)

// NewMiddleware makes every POST safe to retry by default. A POST must carry an
// Idempotency-Key; a repeat with the same key replays the original response.
// Routes can opt out by listing them as "METHOD /path" in optOut, which should
// be rare and deliberate.
func NewMiddleware(db *gorm.DB, optOut ...string) func(http.Handler) http.Handler {
	skip := make(map[string]bool, len(optOut))
	for _, route := range optOut {
		skip[route] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost || skip[r.Method+" "+r.URL.Path] {
				next.ServeHTTP(w, r)
				return
			}

			key := r.Header.Get(HeaderName)
			if key == "" || len(key) > maxKeyLength {
				apijson.Write(w, http.StatusBadRequest, map[string]string{"error": "Idempotency-Key header is required (1-64 characters)"})
				return
			}

			userID, err := auth.UserIDFromContext(r.Context())
			if err != nil {
				apijson.Write(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}

			body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxBodyBytes))
			if err != nil {
				apijson.Write(w, http.StatusBadRequest, map[string]string{"error": "request body could not be read"})
				return
			}
			r.Body = io.NopCloser(bytes.NewReader(body))

			ctx := r.Context()
			hash := Fingerprint(r.Method, r.URL.Path, body)
			s := store{db: db}

			inserted, row, err := s.begin(ctx, principalScope(userID), key, hash)
			if err != nil {
				apijson.Write(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}

			if !inserted {
				switch Decide(row, hash) {
				case ActionMismatch:
					apijson.Write(w, http.StatusUnprocessableEntity, map[string]string{"error": "Idempotency-Key was already used for a different request"})
					return
				case ActionReplay:
					replay(w, row)
					return
				case ActionInProgress:
					// Still running, or it crashed mid-request. Either way this key
					// is not retried; the client must use a fresh key.
					apijson.Write(w, http.StatusConflict, map[string]string{"error": "a request with this Idempotency-Key is still in progress"})
					return
				}
			}

			runAndRecord(w, r, next, s, row.ID)
		})
	}
}

// principalScope names whose keys these are. Today the only principal is a user,
// but the scope is a string so an API client or tenant can be added without a
// schema change.
func principalScope(userID uint) string {
	return "user:" + strconv.FormatUint(uint64(userID), 10)
}

// runAndRecord runs the real handler and stores its response. A 5xx or a panic
// releases the key instead, so a retry runs the handler again.
func runAndRecord(w http.ResponseWriter, r *http.Request, next http.Handler, s store, id uint) {
	rec := newRecorder(w)
	completed := false
	defer func() {
		if completed {
			return
		}
		// Panicking or a 5xx: forget the key so the client can retry.
		_ = s.release(r.Context(), id)
	}()

	next.ServeHTTP(rec, r)
	if rec.status >= 500 {
		return
	}
	if err := s.complete(r.Context(), id, rec.status, rec.contentType(), rec.body.Bytes()); err != nil {
		return
	}
	completed = true
}

func replay(w http.ResponseWriter, row models.IdempotencyKey) {
	if row.ResponseContentType != "" {
		w.Header().Set("Content-Type", row.ResponseContentType)
	}
	w.Header().Set(replayedHeader, strconv.FormatBool(true))
	w.WriteHeader(row.ResponseStatus)
	_, _ = w.Write(row.ResponseBody)
}
