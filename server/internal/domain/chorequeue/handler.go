package chorequeue

import (
	"net/http"
	"time"

	apijson "do-your-dailies/server/internal/api/json"
	"do-your-dailies/server/internal/domain/choreinqueuecompletion"

	"github.com/go-chi/chi/v5"
)

const defaultMaxChores = 5

type Handler struct {
	Now                         func() time.Time
	Store                       Store
	ChoreInQueueCompletionStore choreinqueuecompletion.Store
}

func NewHandler(store Store, choreInQueueCompletionStore choreinqueuecompletion.Store) Handler {
	return Handler{Now: time.Now, Store: store, ChoreInQueueCompletionStore: choreInQueueCompletionStore}
}

func (handler Handler) list(w http.ResponseWriter, r *http.Request) {

	requiredChoresCompletedToday, err := handler.ChoreInQueueCompletionStore.ListBetween(
		time.Now().Add(-24*time.Hour),
		time.Now(),
	)

	if err != nil {
		http.Error(w, "Error listing chore completions in queue", http.StatusInternalServerError)
		return
	}
	choresRemainingToComplete := defaultMaxChores - len(requiredChoresCompletedToday)
	queue, err := handler.Store.ListForCapacityFirstUser(choresRemainingToComplete)
	if err != nil {
		http.Error(w, "Error listing chores in queue", http.StatusInternalServerError)
		return
	}

	apijson.Write(w, http.StatusOK, toAPIChoresInQueue(queue))
}

func (handler Handler) RegisterRoutes(router chi.Router) {
	router.Get("/", handler.list)
}
