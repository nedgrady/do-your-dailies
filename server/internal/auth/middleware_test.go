package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func userIDFromRequest(t *testing.T, simulationModeEnabled bool, header string) uint {
	t.Helper()

	var captured uint
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := UserIDFromContext(r.Context())
		require.NoError(t, err)
		captured = userID
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if header != "" {
		req.Header.Set(SimulatedUserIDHeader, header)
	}
	rr := httptest.NewRecorder()

	NewMiddleware(simulationModeEnabled)(next).ServeHTTP(rr, req)
	return captured
}

func TestMiddlewareIgnoresImpersonationHeaderWhenSimulationModeDisabled(t *testing.T) {
	t.Parallel()
	assert.Equal(t, currentSingleUserID, userIDFromRequest(t, false, "42"))
}

func TestMiddlewareHonorsImpersonationHeaderWhenSimulationModeEnabled(t *testing.T) {
	t.Parallel()
	assert.Equal(t, uint(42), userIDFromRequest(t, true, "42"))
}

func TestMiddlewareFallsBackToDefaultUserWhenHeaderIsMissing(t *testing.T) {
	t.Parallel()
	assert.Equal(t, currentSingleUserID, userIDFromRequest(t, true, ""))
}

func TestMiddlewareFallsBackToDefaultUserWhenHeaderIsMalformed(t *testing.T) {
	t.Parallel()
	assert.Equal(t, currentSingleUserID, userIDFromRequest(t, true, "not-a-number"))
}
