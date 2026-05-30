package requestcontext

import "context"

type userIDKey struct{}

func WithUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

func UserID(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(userIDKey{}).(int)
	return userID, ok
}
