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

	// TODO: How to make sure this is the same timezone as the user's timezone?
	// (when the chore queue is loaded, the FE makes a request in user timezone for 'todays' completions)
	// we want to find 'ttoday's' completions in the chore queue, but we are using UTC time here
	// so we need to unify these timezones somehow
	// add a 'Today' concept?
	previousMidnight := time.Date(handler.Now().Year(), handler.Now().Month(), handler.Now().Day(), 0, 0, 0, 0, time.UTC)
	nextMidnight := previousMidnight.Add(24 * time.Hour)
	requiredChoresCompletedToday, err := handler.ChoreInQueueCompletionStore.ListBetween(
		previousMidnight,
		nextMidnight,
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
