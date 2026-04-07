package api

import (
	"net/http"
	"os"
	"path/filepath"
)

var openAPISpecPathFunc = openAPISpecPath

func (app *Application) openAPISpec(w http.ResponseWriter, r *http.Request) {
	content, err := os.ReadFile(openAPISpecPathFunc())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/yaml")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(content)
}

func (app *Application) swaggerUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = function () {
      window.ui = SwaggerUIBundle({
        url: '/openapi.yaml',
        dom_id: '#swagger-ui'
      });
    };
  </script>
</body>
</html>`))
}

func openAPISpecPath() string {
	candidates := []string{
		filepath.Join("openapi", "openapi.yaml"),
		filepath.Join("server", "openapi", "openapi.yaml"),
		filepath.Join("..", "..", "openapi", "openapi.yaml"),
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return filepath.Clean(candidate)
		}
	}

	return filepath.Clean(candidates[0])
}
