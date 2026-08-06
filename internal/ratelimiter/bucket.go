package ratelimiter

import (
	"sync"
	"time"
)

/*
bucket represents a token bucket for a single client.
Tokens are replenished lazily whenever a request arrives.
*/
type bucket struct {
	mu         sync.Mutex
	tokens     float64
	capacity   float64
	refillRate float64
	lastRefill time.Time
}

/*
newBucket creates a new token bucket.
The bucket starts full, allowing an initial burst of requests.
*/
func newBucket(capacity, refillRate float64) *bucket {
	return &bucket{
		tokens:     capacity,
		capacity:   capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

/*
Allow reports whether a request is permitted.

Before checking for available tokens, the bucket is lazily
refilled based on the time elapsed since the last refill.
*/
func (b *bucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.refill(time.Now())

	if b.tokens < 1 {
		return false
	}

	b.tokens--

	return true
}

/*
refill replenishes tokens based on the elapsed time since
the previous refill.
*/
func (b *bucket) refill(now time.Time) {
	elapsed := now.Sub(b.lastRefill).Seconds()

	b.tokens += elapsed * b.refillRate

	if b.tokens > b.capacity {
		b.tokens = b.capacity
	}

	b.lastRefill = now
}
