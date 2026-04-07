package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	dbmodels "do-your-dailies/server/internal/models"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

var (
	errBadRequest = errors.New("bad request")
	errNotFound   = errors.New("not found")
	errInternal   = errors.New("internal server error")
)

type categorizedError struct {
	category error
	cause    error
}

func (e categorizedError) Error() string {
	if e.cause != nil {
		return e.cause.Error()
	}
	if e.category != nil {
		return e.category.Error()
	}
	return ""
}

func (e categorizedError) Unwrap() error {
	return e.cause
}

func (e categorizedError) Is(target error) bool {
	return target != nil && e.category == target
}

func (app *Application) listChores(w http.ResponseWriter, r *http.Request) {
	chores, err := app.ChoreStore.List()
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, toAPIChores(chores))
}

func (app *Application) createChore(w http.ResponseWriter, r *http.Request) {
	var req CreateChoreRequest
	if err := decodeJSONBody(r, &req); err != nil {
		writeAPIError(w, err)
		return
	}

	chore, err := app.ChoreStore.Create(dbmodels.CreateChoreRequest{
		Name:          req.Name,
		CadenceInDays: req.CadenceInDays,
	})
	if err != nil {
		writeAPIError(w, categorizedError{category: errInternal, cause: err})
		return
	}

	writeJSON(w, http.StatusCreated, toAPIChore(chore))
}

func (app *Application) getChore(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	chore, err := app.ChoreStore.Get(uint(id))
	if err != nil {
		writeAPIError(w, mapStoreError(err))
		return
	}

	writeJSON(w, http.StatusOK, toAPIChore(chore))
}

func (app *Application) updateChore(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	var req UpdateChoreRequest
	if err := decodeJSONBody(r, &req); err != nil {
		writeAPIError(w, err)
		return
	}

	chore, err := app.ChoreStore.Update(uint(id), dbmodels.UpdateChoreRequest{
		Name:          req.Name,
		CadenceInDays: req.CadenceInDays,
	})
	if err != nil {
		writeAPIError(w, mapStoreError(err))
		return
	}

	writeJSON(w, http.StatusOK, toAPIChore(chore))
}

func (app *Application) deleteChore(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if err := app.ChoreStore.Delete(uint(id)); err != nil {
		writeAPIError(w, mapStoreError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func decodeJSONBody(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return categorizedError{category: errBadRequest, cause: err}
	}

	return nil
}

func mapStoreError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return categorizedError{category: errNotFound, cause: err}
	}

	return categorizedError{category: errInternal, cause: err}
}

func writeAPIError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, errBadRequest):
		http.Error(w, errBadRequest.Error(), http.StatusUnprocessableEntity)
	case errors.Is(err, errNotFound):
		http.Error(w, errNotFound.Error(), http.StatusNotFound)
	default:
		http.Error(w, errInternal.Error(), http.StatusInternalServerError)
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func toAPIChores(chores []dbmodels.Chore) []Chore {
	result := make([]Chore, 0, len(chores))
	for _, chore := range chores {
		result = append(result, toAPIChore(chore))
	}
	return result
}

func toAPIChore(chore dbmodels.Chore) Chore {
	return Chore{
		Id:            uint64(chore.ID),
		Name:          chore.Name,
		CadenceInDays: chore.CadenceInDays,
		CreatedAt:     chore.CreatedAt,
		UpdatedAt:     chore.UpdatedAt,
	}
}
