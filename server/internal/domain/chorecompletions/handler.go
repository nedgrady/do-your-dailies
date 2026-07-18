package chorecompletions

import (
	"net/http"
	"slices"
	"time"

	apijson "do-your-dailies/server/internal/api/json"
	"do-your-dailies/server/internal/contracts"
	"do-your-dailies/server/internal/domain/choreinqueuecompletion"
	chorequeue "do-your-dailies/server/internal/domain/chorequeue"
	"do-your-dailies/server/internal/domain/chores"
	apperrors "do-your-dailies/server/internal/errors"

	"gorm.io/gorm"
)

type Handler struct {
	ChoreStore                  chores.Store
	ChoreQueueStore             chorequeue.Store
	ChoreInQueueCompletionStore choreinqueuecompletion.Store
	Now                         func() time.Time
	Store                       Store
	Db                          *gorm.DB
}

func NewHandler(choreStore chores.Store, store Store, choreQueueStore chorequeue.Store, choreInQueueCompletionStore choreinqueuecompletion.Store, db *gorm.DB) Handler {
	return Handler{ChoreStore: choreStore, Now: time.Now, Store: store, ChoreQueueStore: choreQueueStore, ChoreInQueueCompletionStore: choreInQueueCompletionStore, Db: db}
}

func (handler Handler) list(w http.ResponseWriter, r *http.Request) {
	startValue := r.URL.Query().Get("start")
	endValue := r.URL.Query().Get("end")

	if startValue == "" || endValue == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	startTime, err := parseUTCDateTime(startValue)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	endTime, err := parseUTCDateTime(endValue)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if endTime.Before(startTime) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	choreCompletions, listErr := handler.Store.ListByRange(startTime, endTime)
	if listErr != nil {
		apperrors.Write(w, apperrors.Internal(listErr))
		return
	}

	apijson.Write(w, http.StatusOK, toAPIChoreCompletions(choreCompletions))
}

func (handler Handler) create(w http.ResponseWriter, r *http.Request) {
	var req contracts.CreateChoreCompletionRequest
	if err := apijson.DecodeBody(r, &req); err != nil {
		apperrors.Write(w, err)
		return
	}

	if _, err := handler.ChoreStore.Get(uint(req.ChoreId)); err != nil {
		apperrors.Write(w, apperrors.MapStoreError(err))
		return
	}

	choreQueue, err := handler.ChoreQueueStore.ListForCapacityFirstUser(5)
	if err != nil {
		apperrors.Write(w, apperrors.Internal(err))
		return
	}

	isChoreInQueue := slices.ContainsFunc(choreQueue, func(choreInQueue chorequeue.ChoreInQueue) bool {
		return choreInQueue.ChoreID == uint(req.ChoreId)
	})

	if isChoreInQueue {
		handler.ChoreInQueueCompletionStore.Create(choreinqueuecompletion.CreateRequest{ChoreCompletionID: uint(req.ChoreId)})
	}

	choreCompletion, err := handler.Store.Create(CreateRequest{ChoreID: uint(req.ChoreId)})
	if err != nil {
		apperrors.Write(w, apperrors.Internal(err))
		return
	}

	apijson.Write(w, http.StatusCreated, toAPIChoreCompletion(choreCompletion))
}

func parseUTCDateTime(value string) (time.Time, error) {
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return time.Time{}, err
	}

	return parsed.UTC(), nil
}

func startOfDayUTC(value time.Time) time.Time {
	utc := value.UTC()
	return time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC)
}
