package apidocs

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

var openAPISpecPathFunc = openAPISpecPath

func RegisterRoutes(router chi.Router) {
	router.Get("/openapi.yaml", openAPISpec)
}

func openAPISpec(w http.ResponseWriter, r *http.Request) {
	content, err := os.ReadFile(openAPISpecPathFunc())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/yaml")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(content)
}

func openAPISpecPath() string {
	candidates := []string{
		filepath.Join("openapi", "openapi.yaml"),
		filepath.Join("server", "openapi", "openapi.yaml"),
		filepath.Join("..", "..", "openapi", "openapi.yaml"),
		filepath.Join("..", "..", "..", "openapi", "openapi.yaml"),
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return filepath.Clean(candidate)
		}
	}

	return filepath.Clean(candidates[0])
}
