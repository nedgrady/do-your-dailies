package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *Application) listChores(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"message": "list chores"})
}

func (app *Application) createChore(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusCreated, map[string]string{"message": "create chore"})
}

func (app *Application) getChore(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	writeJSON(w, http.StatusOK, map[string]string{"id": id})
}

func writeJSON(w http.ResponseWriter, status int, payload map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
