package chorecompletions

import (
	"context"
	"errors"
	"slices"
	"time"

	"do-your-dailies/server/internal/contracts"
	"do-your-dailies/server/internal/domain/choreinqueuecompletion"
	chorequeue "do-your-dailies/server/internal/domain/chorequeue"
	"do-your-dailies/server/internal/domain/chores"

	"gorm.io/gorm"
)

type ChoreCompletionHandler struct {
	ChoreStore                  chores.Store
	ChoreQueueStore             chorequeue.Store
	ChoreInQueueCompletionStore choreinqueuecompletion.Store
	Now                         func() time.Time
	Store                       Store
	Db                          *gorm.DB
}

func NewHandler(choreStore chores.Store, store Store, choreQueueStore chorequeue.Store, choreInQueueCompletionStore choreinqueuecompletion.Store, db *gorm.DB) ChoreCompletionHandler {
	return ChoreCompletionHandler{ChoreStore: choreStore, Now: time.Now, Store: store, ChoreQueueStore: choreQueueStore, ChoreInQueueCompletionStore: choreInQueueCompletionStore, Db: db}
}

func (handler ChoreCompletionHandler) ListChoreCompletions(ctx context.Context, request contracts.ListChoreCompletionsRequestObject) (contracts.ListChoreCompletionsResponseObject, error) {
	if request.Params.End.Before(request.Params.Start.Time) {
		return contracts.ListChoreCompletions400Response{}, nil
	}

	choreCompletions, err := handler.Store.ListByRange(request.Params.Start.Time, request.Params.End.Time)
	if err != nil {
		return nil, err
	}

	return contracts.ListChoreCompletions200JSONResponse(toAPIChoreCompletions(choreCompletions)), nil
}

func (handler ChoreCompletionHandler) CreateChoreCompletion(ctx context.Context, request contracts.CreateChoreCompletionRequestObject) (contracts.CreateChoreCompletionResponseObject, error) {
	choreId := uint(request.Body.ChoreId)

	if _, err := handler.ChoreStore.Get(choreId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return contracts.CreateChoreCompletion404Response{}, nil
		}
		return nil, err
	}

	choreQueue, err := handler.ChoreQueueStore.ListForCapacityFirstUser(5)
	if err != nil {
		return nil, err
	}

	isChoreInQueue := slices.ContainsFunc(choreQueue, func(choreInQueue chorequeue.ChoreInQueue) bool {
		return choreInQueue.ChoreID == choreId
	})

	if isChoreInQueue {
		handler.ChoreInQueueCompletionStore.Create(choreinqueuecompletion.CreateRequest{ChoreCompletionID: choreId})
	}

	choreCompletion, err := handler.Store.Create(CreateRequest{ChoreID: choreId})
	if err != nil {
		return nil, err
	}

	return contracts.CreateChoreCompletion201JSONResponse(toAPIChoreCompletion(choreCompletion)), nil
}
