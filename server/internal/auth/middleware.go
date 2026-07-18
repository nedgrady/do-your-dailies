package auth

import "net/http"

// TODO: replace with real auth (session/token lookup) once there's more than one user.
const currentSingleUserID uint = 1

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := WithUserID(r.Context(), currentSingleUserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
