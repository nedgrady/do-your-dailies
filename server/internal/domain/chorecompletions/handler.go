package chorecompletions

import (
	"net/http"
	"slices"
	"time"

	apijson "do-your-dailies/server/internal/api/json"
	"do-your-dailies/server/internal/contracts"
	chorequeue "do-your-dailies/server/internal/domain/chorequeue"
	"do-your-dailies/server/internal/domain/chores"
	apperrors "do-your-dailies/server/internal/errors"
)

type Handler struct {
	ChoreStore      chores.Store
	ChoreQueueStore chorequeue.Store
	Now             func() time.Time
	Store           Store
}

func NewHandler(choreStore chores.Store, store Store) Handler {
	return Handler{ChoreStore: choreStore, Now: time.Now, Store: store}
}

func (handler Handler) list(w http.ResponseWriter, r *http.Request) {
	day, err := optionalDateQuery(r, "date", startOfDayUTC(handler.Now()))
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	choreCompletions, listErr := handler.Store.ListByDay(day)
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

	if isChoreInQueue  {
		
	}

	choreCompletion, err := handler.Store.Create(CreateRequest{ChoreID: uint(req.ChoreId)})
	if err != nil {
		apperrors.Write(w, apperrors.Internal(err))
		return
	}

	apijson.Write(w, http.StatusCreated, toAPIChoreCompletion(choreCompletion))
}

func optionalDateQuery(r *http.Request, key string, defaultValue time.Time) (time.Time, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return startOfDayUTC(defaultValue), nil
	}

	date, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, err
	}

	return startOfDayUTC(date), nil
}

func startOfDayUTC(value time.Time) time.Time {
	utc := value.UTC()
	return time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC)
}
