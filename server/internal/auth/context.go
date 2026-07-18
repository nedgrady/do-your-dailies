package auth

import (
	"context"
	"errors"
)

type contextKey string

const userIDContextKey contextKey = "userID"

func WithUserID(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

func UserIDFromContext(ctx context.Context) (uint, error) {
	userID, ok := ctx.Value(userIDContextKey).(uint)
	if !ok {
		return 0, errors.New("no user id in context")
	}
	return userID, nil
}
