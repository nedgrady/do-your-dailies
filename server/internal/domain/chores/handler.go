package chores

import (
	"context"
	"errors"

	"do-your-dailies/server/internal/auth"
	"do-your-dailies/server/internal/contracts"

	"gorm.io/gorm"
)

type ChoreHandler struct {
	Store Store
}

func NewHandler(store Store) ChoreHandler {
	return ChoreHandler{Store: store}
}

func (handler ChoreHandler) ListChores(ctx context.Context, request contracts.ListChoresRequestObject) (contracts.ListChoresResponseObject, error) {
	userID, err := auth.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	chores, err := handler.Store.List(ctx, userID)
	if err != nil {
		return nil, err
	}

	return contracts.ListChores200JSONResponse(ToAPIChores(chores)), nil
}

func (handler ChoreHandler) CreateChore(ctx context.Context, request contracts.CreateChoreRequestObject) (contracts.CreateChoreResponseObject, error) {
	userID, err := auth.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := validateCadence(request.Body.CadenceInDays); err != nil {
		return contracts.CreateChore422Response{}, nil
	}

	chore, err := handler.Store.Create(ctx, userID, CreateRequest{
		Name:          request.Body.Name,
		CadenceInDays: request.Body.CadenceInDays,
	})
	if err != nil {
		return nil, err
	}

	return contracts.CreateChore201JSONResponse(ToAPIChore(chore)), nil
}

func (handler ChoreHandler) BulkCreateChores(ctx context.Context, request contracts.BulkCreateChoresRequestObject) (contracts.BulkCreateChoresResponseObject, error) {
	userID, err := auth.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var rowErrors []contracts.BulkCreateChoresRowError
	requests := make([]CreateRequest, 0, len(request.Body.Chores))
	for i, chore := range request.Body.Chores {
		if err := validateCadence(chore.CadenceInDays); err != nil {
			rowErrors = append(rowErrors, contracts.BulkCreateChoresRowError{Index: i, Message: err.Error()})
			continue
		}
		if chore.Name == "" {
			rowErrors = append(rowErrors, contracts.BulkCreateChoresRowError{Index: i, Message: "name must not be empty"})
			continue
		}
		requests = append(requests, CreateRequest{Name: chore.Name, CadenceInDays: chore.CadenceInDays})
	}

	if len(rowErrors) > 0 {
		return contracts.BulkCreateChores422JSONResponse{Errors: rowErrors}, nil
	}

	created, err := handler.Store.CreateMany(ctx, userID, requests)
	if err != nil {
		return nil, err
	}

	return contracts.BulkCreateChores201JSONResponse(ToAPIChores(created)), nil
}

func (handler ChoreHandler) GetChore(ctx context.Context, request contracts.GetChoreRequestObject) (contracts.GetChoreResponseObject, error) {
	userID, err := auth.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	chore, err := handler.Store.Get(ctx, userID, uint(request.Id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return contracts.GetChore404Response{}, nil
		}
		return nil, err
	}

	return contracts.GetChore200JSONResponse(ToAPIChore(chore)), nil
}

func (handler ChoreHandler) UpdateChore(ctx context.Context, request contracts.UpdateChoreRequestObject) (contracts.UpdateChoreResponseObject, error) {
	userID, err := auth.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if request.Body.CadenceInDays != nil {
		if err := validateCadence(*request.Body.CadenceInDays); err != nil {
			return contracts.UpdateChore422Response{}, nil
		}
	}

	chore, err := handler.Store.Update(ctx, userID, uint(request.Id), UpdateRequest{
		Name:          request.Body.Name,
		CadenceInDays: request.Body.CadenceInDays,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return contracts.UpdateChore404Response{}, nil
		}
		return nil, err
	}

	return contracts.UpdateChore200JSONResponse(ToAPIChore(chore)), nil
}

func (handler ChoreHandler) DeleteChore(ctx context.Context, request contracts.DeleteChoreRequestObject) (contracts.DeleteChoreResponseObject, error) {
	userID, err := auth.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := handler.Store.Delete(ctx, userID, uint(request.Id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return contracts.DeleteChore404Response{}, nil
		}
		return nil, err
	}

	return contracts.DeleteChore204Response{}, nil
}

func validateCadence(cadenceInDays int) error {
	if cadenceInDays <= 0 {
		return errors.New("cadenceInDays must be greater than zero")
	}

	return nil
}
