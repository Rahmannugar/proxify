package ratelimiter

// Limiter applies token bucket rate limiting for clients.
type Limiter struct {
	store *Store
}

// New creates a new rate limiter.
func New(requestsPerSecond, burst int) *Limiter {
	return &Limiter{
		store: NewStore(
			float64(burst),
			float64(requestsPerSecond),
		),
	}
}

/*
Allow reports whether the given client is permitted
to make another request.
*/
func (l *Limiter) Allow(clientID string) bool {
	return l.store.Bucket(clientID).Allow()
}
