package api

import (
	"net/http"

	"do-your-dailies/server/internal/docs"
	"do-your-dailies/server/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"
)

type Application struct {
	Router               *chi.Mux
	DB                   *gorm.DB
	ChoreStore           store.ChoreStore
	ChoreCompletionStore store.ChoreCompletionStore
}

func New(db *gorm.DB) *Application {
	app := &Application{DB: db}
	if db != nil {
		app.ChoreStore = store.NewGormChoreStore(db)
		app.ChoreCompletionStore = store.NewGormChoreCompletionStore(db)
	}
	app.Router = app.setupRoutes()
	return app
}

func (app *Application) setupRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", app.healthCheck)
	docs.NewHandler().RegisterRoutes(r)

	r.Route("/api", func(r chi.Router) {
		r.Route("/chores", func(r chi.Router) {
			r.Get("/", app.listChores)
			r.Post("/", app.createChore)
			r.Get("/{id}", app.getChore)
			r.Put("/{id}", app.updateChore)
			r.Delete("/{id}", app.deleteChore)
		})

		r.Route("/chore-completions", func(r chi.Router) {
			r.Post("/", app.createChoreCompletion)
		})
	})

	return r
}

func (app *Application) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("healthy"))
}
