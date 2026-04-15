package chorequeue

import (
	"net/http"
	"strconv"
	"time"

	apijson "do-your-dailies/server/internal/api/json"
	"do-your-dailies/server/internal/errors"

	"github.com/go-chi/chi/v5"
)

const defaultMaxChores = 10

type Handler struct {
	Now   func() time.Time
	Store Store
}

func NewHandler(store Store) Handler {
	return Handler{Now: time.Now, Store: store}
}

func (handler Handler) list(w http.ResponseWriter, r *http.Request) {
	dayOffset, err := optionalIntQuery(r, "dayOffset", 0)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	maxChores, err := optionalIntQuery(r, "maxChores", defaultMaxChores)
	if err != nil || maxChores <= 0 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	targetDay := startOfDayUTC(handler.Now()).AddDate(0, 0, dayOffset)
	queue, listErr := handler.Store.List(targetDay, maxChores)
	if listErr != nil {
		errors.Write(w, errors.MapStoreError(listErr))
		return
	}

	apijson.Write(w, http.StatusOK, toAPIChores(queue))
}

func optionalIntQuery(r *http.Request, key string, defaultValue int) (int, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue, nil
	}

	return strconv.Atoi(value)
}

func (handler Handler) RegisterRoutes(router chi.Router) {
	router.Get("/", handler.list)
}
