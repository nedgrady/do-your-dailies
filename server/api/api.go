package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Application struct {
	Router *chi.Mux
}

func New() *Application {
	app := &Application{}
	app.Router = app.setupRoutes()
	return app
}

func (app *Application) setupRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Route("/api", func(r chi.Router) {
		r.Route("/chores", func(r chi.Router) {
			r.Get("/", app.listChores)
			r.Post("/", app.createChore)
			r.Get("/{id}", app.getChore)
		})
	})

	return r
}
