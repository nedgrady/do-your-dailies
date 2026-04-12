package chores

import "github.com/go-chi/chi/v5"

func (handler Handler) RegisterRoutes(router chi.Router) {
	router.Get("/", handler.list)
	router.Post("/", handler.create)
	router.Get("/{id}", handler.get)
	router.Put("/{id}", handler.update)
	router.Delete("/{id}", handler.delete)
}
