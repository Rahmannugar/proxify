package circuitbreaker

// State represents the current circuit breaker state.
type State int

const (
	Closed State = iota
	Open
	HalfOpen
)
