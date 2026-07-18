package api

import (
	"net/http"

	apijson "do-your-dailies/server/internal/api/json"
	"do-your-dailies/server/internal/apidocs"
	"do-your-dailies/server/internal/auth"
	"do-your-dailies/server/internal/contracts"
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
	router.Use(corsMiddleware)
	router.Use(middleware.Recoverer)
	router.Use(auth.Middleware)

	apidocs.RegisterRoutes(router)
	swagger.RegisterRoutes(router)

	server := Server{
		ChoreHandler: chores.NewHandler(app.ChoreStore),
		ChoreCompletionHandler: chorecompletions.NewHandler(
			app.ChoreStore,
			app.ChoreCompletionStore,
			app.ChoreQueueStore,
			app.ChoreInQueueCompletionStore,
			app.Db,
		),
		ChoreQueueHandler: chorequeue.NewHandler(app.ChoreQueueStore, app.ChoreInQueueCompletionStore),
	}

	strictHandler := contracts.NewStrictHandlerWithOptions(server, nil, contracts.StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			apijson.Write(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			apijson.Write(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		},
	})

	contracts.HandlerFromMux(strictHandler, router)

	return router
}
