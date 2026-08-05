package ratelimiter

import "sync"

// Store manages token buckets for all clients.
type Store struct {
	mu sync.RWMutex

	buckets map[string]*bucket

	capacity   float64
	refillRate float64
}

// NewStore creates a new bucket store.
func NewStore(capacity, refillRate float64) *Store {
	return &Store{
		buckets:    make(map[string]*bucket),
		capacity:   capacity,
		refillRate: refillRate,
	}
}

/*
Bucket returns the token bucket for the given client.
A new bucket is created if one does not already exist.
*/
func (s *Store) Bucket(clientID string) *bucket {
	s.mu.RLock()

	bucket, ok := s.buckets[clientID]

	s.mu.RUnlock()

	if ok {
		return bucket
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	/*
		Another goroutine may have created the bucket while
		we were waiting for the write lock.
	*/
	bucket, ok = s.buckets[clientID]
	if ok {
		return bucket
	}

	bucket = newBucket(
		s.capacity,
		s.refillRate,
	)

	s.buckets[clientID] = bucket

	return bucket
}
