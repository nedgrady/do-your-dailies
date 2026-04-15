package api

import (
	"do-your-dailies/server/internal/apidocs"
	"do-your-dailies/server/internal/domain/chorecompletions"
	"do-your-dailies/server/internal/domain/chorequeue"
	"do-your-dailies/server/internal/domain/chores"
	"do-your-dailies/server/internal/logging"
	"do-your-dailies/server/internal/swagger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Application) setupRoutes() *chi.Mux {
	router := chi.NewRouter()
	logging.Apply(router)
	router.Use(middleware.Recoverer)

	router.Get("/health", healthCheck)
	apidocs.RegisterRoutes(router)
	swagger.RegisterRoutes(router)

	choreHandler := chores.NewHandler(app.ChoreStore)
	choreCompletionHandler := chorecompletions.NewHandler(app.ChoreStore, app.ChoreCompletionStore)
	choreQueueHandler := chorequeue.NewHandler(app.ChoreQueueStore)

	router.Route("/api", func(router chi.Router) {
		router.Route("/chores", choreHandler.RegisterRoutes)
		router.Route("/chore-completions", choreCompletionHandler.RegisterRoutes)
		router.Route("/chore-queue", choreQueueHandler.RegisterRoutes)
	})

	return router
}
