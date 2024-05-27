package api

import "context"

// contextKey represents an internal key for adding context fields.
// This is considered best practice as it prevents other packages from
// interfering with our context keys.
type contextKey int

// List of context keys.
// These are used to store request-scoped information.
const (
	// Stores the current logged in user in the context.
	middlewareDataContextKey = contextKey(iota + 1)
)

// MiddlewareDataFromContext returns the current middleware data.
func MiddlewareDataFromContext(ctx context.Context) *Data {
	data, _ := ctx.Value(middlewareDataContextKey).(*Data)
	return data
}

// NewContextWithMiddlewareData returns a new context with the given middleware data.
func NewContextWithMiddlewareData(ctx context.Context, data *Data) context.Context {
	return context.WithValue(ctx, middlewareDataContextKey, data)
}
