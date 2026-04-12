package chorecompletions

import (
	"net/http"

	apijson "do-your-dailies/server/internal/api/json"
	"do-your-dailies/server/internal/contracts"
	"do-your-dailies/server/internal/domain/chores"
	apperrors "do-your-dailies/server/internal/errors"
)

type Handler struct {
	ChoreStore chores.Store
	Store      Store
}

func NewHandler(choreStore chores.Store, store Store) Handler {
	return Handler{ChoreStore: choreStore, Store: store}
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

	choreCompletion, err := handler.Store.Create(CreateRequest{ChoreID: uint(req.ChoreId)})
	if err != nil {
		apperrors.Write(w, apperrors.Internal(err))
		return
	}

	apijson.Write(w, http.StatusCreated, toAPIChoreCompletion(choreCompletion))
}
