package circuitbreaker

import "errors"

// ErrOpen indicates the circuit breaker is open.
var ErrOpen = errors.New("circuit breaker is open")
