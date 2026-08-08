package chorequeue

import (
	"context"
	"time"

	"do-your-dailies/server/internal/auth"
	"do-your-dailies/server/internal/contracts"
	"do-your-dailies/server/internal/domain/choreinqueuecompletion"
)

const defaultMaxChores = 5

type ChoreQueueHandler struct {
	Now                         func() time.Time
	Store                       Store
	ChoreInQueueCompletionStore choreinqueuecompletion.Store
}

func NewHandler(store Store, choreInQueueCompletionStore choreinqueuecompletion.Store) ChoreQueueHandler {
	return ChoreQueueHandler{Now: time.Now, Store: store, ChoreInQueueCompletionStore: choreInQueueCompletionStore}
}

func (handler ChoreQueueHandler) ListChoreQueue(ctx context.Context, request contracts.ListChoreQueueRequestObject) (contracts.ListChoreQueueResponseObject, error) {
	userID, err := auth.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// TODO: How to make sure this is the same timezone as the user's timezone?
	// (when the chore queue is loaded, the FE makes a request in user timezone for 'todays' completions)
	// we want to find 'ttoday's' completions in the chore queue, but we are using UTC time here
	// so we need to unify these timezones somehow
	// add a 'Today' concept?
	previousMidnight := time.Date(handler.Now().Year(), handler.Now().Month(), handler.Now().Day(), 0, 0, 0, 0, time.UTC)
	nextMidnight := previousMidnight.Add(24 * time.Hour)
	requiredChoresCompletedToday, err := handler.ChoreInQueueCompletionStore.ListBetween(
		ctx,
		userID,
		previousMidnight,
		nextMidnight,
	)
	if err != nil {
		return nil, err
	}

	choresRemainingToComplete := defaultMaxChores - len(requiredChoresCompletedToday)
	queue, err := handler.Store.ListForCapacityFirstUser(ctx, userID, choresRemainingToComplete)
	if err != nil {
		return nil, err
	}

	return contracts.ListChoreQueue200JSONResponse(toAPIChoresInQueue(queue)), nil
}
