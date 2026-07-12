package chorequeue

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	apijson "do-your-dailies/server/internal/api/json"

	"github.com/go-chi/chi/v5"
)

const defaultMaxChores = 5

type Handler struct {
	Now   func() time.Time
	Store Store
}

func NewHandler(store Store) Handler {
	return Handler{Now: time.Now, Store: store}
}

func (handler Handler) list(w http.ResponseWriter, r *http.Request) {
	maxChores, err := parseMaxChores(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	queue, err := handler.Store.ListForCapacityFirstUser(maxChores)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	apijson.Write(w, http.StatusOK, toAPIChoresInQueue(queue))
}

func parseMaxChores(r *http.Request) (int, error) {
	raw := r.URL.Query().Get("maxChores")
	if raw == "" {
		return defaultMaxChores, nil
	}

	maxChores, err := strconv.Atoi(raw)
	if err != nil || maxChores <= 0 {
		return 0, errors.New("invalid maxChores")
	}

	return maxChores, nil
}

func (handler Handler) RegisterRoutes(router chi.Router) {
	router.Get("/", handler.list)
}
