package chorecompletions

import "github.com/go-chi/chi/v5"

func (handler Handler) RegisterRoutes(router chi.Router) {
	router.Post("/", handler.create)
}
