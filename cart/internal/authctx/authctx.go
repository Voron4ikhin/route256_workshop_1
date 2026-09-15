package authctx

import "context"

type contextKey struct{}

var tokenKey = contextKey{}

func WithToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenKey, token)
}

func Token(ctx context.Context) string {
	token, _ := ctx.Value(tokenKey).(string)
	return token
}
