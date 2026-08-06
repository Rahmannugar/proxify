package circuitbreaker

import (
	"sync"
	"time"
)

// Breaker protects an upstream from repeated failures.
type Breaker struct {
	mu sync.Mutex

	state State

	failureThreshold int
	successThreshold int
	recoveryTimeout  time.Duration

	failures  int
	successes int

	lastFailure time.Time
}

// New creates a new circuit breaker.
func New(
	failureThreshold int,
	successThreshold int,
	recoveryTimeout time.Duration,
) *Breaker {
	return &Breaker{
		state:            Closed,
		failureThreshold: failureThreshold,
		successThreshold: successThreshold,
		recoveryTimeout:  recoveryTimeout,
	}
}

// Allow reports whether a request may proceed.
func (b *Breaker) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case Closed:
		return true

	case Open:
		if time.Since(b.lastFailure) >= b.recoveryTimeout {
			b.state = HalfOpen
			b.successes = 0

			return true
		}

		return false

	case HalfOpen:
		return true

	default:
		return false
	}
}

// Success records a successful request.
func (b *Breaker) Success() {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case Closed:
		b.failures = 0

	case HalfOpen:
		b.successes++

		if b.successes >= b.successThreshold {
			b.state = Closed
			b.failures = 0
			b.successes = 0
		}
	}
}

// Failure records a failed request.
func (b *Breaker) Failure() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.lastFailure = time.Now()

	switch b.state {
	case Closed:
		b.failures++

		if b.failures >= b.failureThreshold {
			b.state = Open
		}

	case HalfOpen:
		b.state = Open
		b.successes = 0
	}
}

// State returns the current circuit breaker state.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.state
}
