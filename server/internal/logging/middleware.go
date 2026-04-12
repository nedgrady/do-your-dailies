package logging

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Apply(router chi.Router) {
	router.Use(middleware.Logger)
}
