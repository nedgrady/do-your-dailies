package auth

import (
	"net/http"
	"strconv"
)

// TODO: replace with real auth (session/token lookup) once there's more than one user.
const currentSingleUserID uint = 1

// SimulatedUserIDHeader lets a caller act as an arbitrary user id, for
// exploring cold-start scheduling scenarios without a real second account.
// Only honored when simulationModeEnabled is true; otherwise every request
// is stamped with currentSingleUserID exactly as before, regardless of any
// header sent - this is what keeps the mechanism inert by default.
const SimulatedUserIDHeader = "X-Simulated-User-Id"

func NewMiddleware(simulationModeEnabled bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := currentSingleUserID
			if simulationModeEnabled {
				if simulated, err := strconv.ParseUint(r.Header.Get(SimulatedUserIDHeader), 10, 64); err == nil {
					userID = uint(simulated)
				}
			}
			ctx := WithUserID(r.Context(), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
