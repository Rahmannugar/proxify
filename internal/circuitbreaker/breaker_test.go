package circuitbreaker

import (
	"testing"
	"time"
)

func TestBreaker_OpensAfterFailureThreshold(t *testing.T) {
	breaker := New(
		3,
		2,
		time.Second,
	)

	breaker.Failure()
	breaker.Failure()

	if breaker.State() != Closed {
		t.Fatal("expected breaker to remain closed")
	}

	breaker.Failure()

	if breaker.State() != Open {
		t.Fatal("expected breaker to open")
	}
}

func TestBreaker_TransitionsToHalfOpen(t *testing.T) {
	breaker := New(
		1,
		1,
		10*time.Millisecond,
	)

	breaker.Failure()

	time.Sleep(20 * time.Millisecond)

	if !breaker.Allow() {
		t.Fatal("expected breaker to allow request")
	}

	if breaker.State() != HalfOpen {
		t.Fatal("expected breaker to transition to half-open")
	}
}

func TestBreaker_ClosesAfterSuccessfulRecovery(t *testing.T) {
	breaker := New(
		1,
		2,
		10*time.Millisecond,
	)

	breaker.Failure()

	time.Sleep(20 * time.Millisecond)

	breaker.Allow()

	breaker.Success()

	if breaker.State() != HalfOpen {
		t.Fatal("expected breaker to remain half-open")
	}

	breaker.Success()

	if breaker.State() != Closed {
		t.Fatal("expected breaker to close")
	}
}

func TestBreaker_ReopensAfterFailureInHalfOpen(t *testing.T) {
	breaker := New(
		1,
		1,
		10*time.Millisecond,
	)

	breaker.Failure()

	time.Sleep(20 * time.Millisecond)

	breaker.Allow()

	breaker.Failure()

	if breaker.State() != Open {
		t.Fatal("expected breaker to reopen")
	}
}
