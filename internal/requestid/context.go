package requestid

import "context"

type contextKey string

const key contextKey = "request-id"

// FromContext returns the request ID stored in the context.
func FromContext(ctx context.Context) string {
	requestID, _ := ctx.Value(key).(string)

	return requestID
}

// NewContext returns a copy of the context containing the request ID.
func NewContext(
	ctx context.Context,
	requestID string,
) context.Context {
	return context.WithValue(
		ctx,
		key,
		requestID,
	)
}
