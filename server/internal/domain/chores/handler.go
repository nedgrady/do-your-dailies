package chores

import (
	"context"
	"errors"

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
	chores, err := handler.Store.List()
	if err != nil {
		return nil, err
	}

	return contracts.ListChores200JSONResponse(toAPIChores(chores)), nil
}

func (handler ChoreHandler) CreateChore(ctx context.Context, request contracts.CreateChoreRequestObject) (contracts.CreateChoreResponseObject, error) {
	if err := validateCadence(request.Body.CadenceInDays); err != nil {
		return contracts.CreateChore422Response{}, nil
	}

	chore, err := handler.Store.Create(CreateRequest{
		Name:          request.Body.Name,
		CadenceInDays: request.Body.CadenceInDays,
	})
	if err != nil {
		return nil, err
	}

	return contracts.CreateChore201JSONResponse(toAPIChore(chore)), nil
}

func (handler ChoreHandler) GetChore(ctx context.Context, request contracts.GetChoreRequestObject) (contracts.GetChoreResponseObject, error) {
	chore, err := handler.Store.Get(uint(request.Id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return contracts.GetChore404Response{}, nil
		}
		return nil, err
	}

	return contracts.GetChore200JSONResponse(toAPIChore(chore)), nil
}

func (handler ChoreHandler) UpdateChore(ctx context.Context, request contracts.UpdateChoreRequestObject) (contracts.UpdateChoreResponseObject, error) {
	if request.Body.CadenceInDays != nil {
		if err := validateCadence(*request.Body.CadenceInDays); err != nil {
			return contracts.UpdateChore422Response{}, nil
		}
	}

	chore, err := handler.Store.Update(uint(request.Id), UpdateRequest{
		Name:          request.Body.Name,
		CadenceInDays: request.Body.CadenceInDays,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return contracts.UpdateChore404Response{}, nil
		}
		return nil, err
	}

	return contracts.UpdateChore200JSONResponse(toAPIChore(chore)), nil
}

func (handler ChoreHandler) DeleteChore(ctx context.Context, request contracts.DeleteChoreRequestObject) (contracts.DeleteChoreResponseObject, error) {
	if err := handler.Store.Delete(uint(request.Id)); err != nil {
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
