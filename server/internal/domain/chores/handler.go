package chores

import (
	"errors"
	"net/http"
	"strconv"

	apijson "do-your-dailies/server/internal/api/json"
	"do-your-dailies/server/internal/contracts"
	apperrors "do-your-dailies/server/internal/errors"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Store Store
}

func NewHandler(store Store) Handler {
	return Handler{Store: store}
}

func (handler Handler) list(w http.ResponseWriter, r *http.Request) {
	chores, err := handler.Store.List()
	if err != nil {
		apperrors.Write(w, apperrors.MapStoreError(err))
		return
	}

	apijson.Write(w, http.StatusOK, toAPIChores(chores))
}

func (handler Handler) create(w http.ResponseWriter, r *http.Request) {
	var req contracts.CreateChoreRequest
	if err := apijson.DecodeBody(r, &req); err != nil {
		apperrors.Write(w, err)
		return
	}
	if err := validateCadence(req.CadenceInDays); err != nil {
		apperrors.Write(w, apperrors.BadRequest(err))
		return
	}

	chore, err := handler.Store.Create(CreateRequest{
		Name:          req.Name,
		CadenceInDays: req.CadenceInDays,
	})
	if err != nil {
		apperrors.Write(w, apperrors.Internal(err))
		return
	}

	apijson.Write(w, http.StatusCreated, toAPIChore(chore))
}

func (handler Handler) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	chore, err := handler.Store.Get(uint(id))
	if err != nil {
		apperrors.Write(w, apperrors.MapStoreError(err))
		return
	}

	apijson.Write(w, http.StatusOK, toAPIChore(chore))
}

func (handler Handler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	var req contracts.UpdateChoreRequest
	if err := apijson.DecodeBody(r, &req); err != nil {
		apperrors.Write(w, err)
		return
	}
	if req.CadenceInDays != nil {
		if err := validateCadence(*req.CadenceInDays); err != nil {
			apperrors.Write(w, apperrors.BadRequest(err))
			return
		}
	}

	chore, err := handler.Store.Update(uint(id), UpdateRequest{
		Name:          req.Name,
		CadenceInDays: req.CadenceInDays,
	})
	if err != nil {
		apperrors.Write(w, apperrors.MapStoreError(err))
		return
	}

	apijson.Write(w, http.StatusOK, toAPIChore(chore))
}

func (handler Handler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if err := handler.Store.Delete(uint(id)); err != nil {
		apperrors.Write(w, apperrors.MapStoreError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func validateCadence(cadenceInDays int) error {
	if cadenceInDays <= 0 {
		return errors.New("cadenceInDays must be greater than zero")
	}

	return nil
}
