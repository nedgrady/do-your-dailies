package swagger

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func newSwaggerRouter() *chi.Mux {
	router := chi.NewRouter()
	RegisterRoutes(router)
	return router
}

func TestSwaggerRouteRedirects(t *testing.T) {
	t.Parallel()
	router := newSwaggerRouter()

	req := httptest.NewRequest(http.MethodGet, "/swagger", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMovedPermanently, rr.Code)
}

func TestSwaggerRouteRedirectsToTrailingSlash(t *testing.T) {
	t.Parallel()
	router := newSwaggerRouter()

	req := httptest.NewRequest(http.MethodGet, "/swagger", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, "/swagger/", rr.Header().Get("Location"))
}
