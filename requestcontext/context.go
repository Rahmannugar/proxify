package requestcontext

import "context"

/*
contextKey represents a request context key.

Using a dedicated type avoids collisions with keys
defined by other packages.
*/
type contextKey string

const upstreamIDKey contextKey = "upstreamID"

/*
WithUpstreamID stores the selected upstream ID in the request context.
*/
func WithUpstreamID(ctx context.Context, upstreamID string) context.Context {
	return context.WithValue(ctx, upstreamIDKey, upstreamID)
}

/*
UpstreamID retrieves the selected upstream ID from the request context.
*/
func UpstreamID(ctx context.Context) (string, bool) {
	upstreamID, ok := ctx.Value(upstreamIDKey).(string)

	return upstreamID, ok
}
