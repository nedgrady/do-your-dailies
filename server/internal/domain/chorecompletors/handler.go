package chorecompletors

import (
	"context"
	"strings"

	"do-your-dailies/server/internal/auth"
	"do-your-dailies/server/internal/contracts"
)

type ChoreCompletorHandler struct {
	Store Store
}

func NewHandler(store Store) ChoreCompletorHandler {
	return ChoreCompletorHandler{Store: store}
}

func (handler ChoreCompletorHandler) ListChoreCompletors(ctx context.Context, request contracts.ListChoreCompletorsRequestObject) (contracts.ListChoreCompletorsResponseObject, error) {
	userID, err := auth.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	completors, err := handler.Store.List(ctx, userID)
	if err != nil {
		return nil, err
	}

	return contracts.ListChoreCompletors200JSONResponse(ToAPIChoreCompletors(completors)), nil
}

func (handler ChoreCompletorHandler) CreateChoreCompletor(ctx context.Context, request contracts.CreateChoreCompletorRequestObject) (contracts.CreateChoreCompletorResponseObject, error) {
	userID, err := auth.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	name := strings.TrimSpace(request.Body.Name)
	if name == "" {
		return contracts.CreateChoreCompletor422Response{}, nil
	}

	completor, err := handler.Store.Create(ctx, userID, name)
	if err != nil {
		return nil, err
	}

	return contracts.CreateChoreCompletor201JSONResponse(ToAPIChoreCompletor(completor)), nil
}
