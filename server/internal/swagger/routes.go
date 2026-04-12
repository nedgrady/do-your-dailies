package swagger

import "github.com/go-chi/chi/v5"

func RegisterRoutes(router chi.Router) {
	router.Get("/swagger", redirect)
	router.Get("/swagger/", ui)
}
