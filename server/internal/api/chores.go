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
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusUnprocessableEntity)
		return
	}
	chore, err := app.ChoreStore.Create(dbmodels.CreateChoreRequest{
		Name:          req.Name,
		CadenceInDays: req.CadenceInDays,
	})
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
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
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusUnprocessableEntity)
		return
	}
	chore, err := app.ChoreStore.Update(uint(id), dbmodels.UpdateChoreRequest{
		Name:          req.Name,
		CadenceInDays: req.CadenceInDays,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
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
