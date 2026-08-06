package ratelimiter

import (
	"sync"
	"testing"
	"time"
)

func TestLimiter_Allow_AllowsRequestsWhileTokensExist(t *testing.T) {
	limiter := New(1, 2)

	if !limiter.Allow("client-1") {
		t.Fatal("expected first request to be allowed")
	}

	if !limiter.Allow("client-1") {
		t.Fatal("expected second request to be allowed")
	}

	if limiter.Allow("client-1") {
		t.Fatal("expected third request to be rejected")
	}
}

func TestLimiter_Allow_RefillsTokens(t *testing.T) {
	limiter := New(2, 1)

	if !limiter.Allow("client-1") {
		t.Fatal("expected request to be allowed")
	}

	time.Sleep(600 * time.Millisecond)

	if !limiter.Allow("client-1") {
		t.Fatal("expected token to be refilled")
	}
}

func TestLimiter_Allow_IsPerClient(t *testing.T) {
	limiter := New(1, 1)

	if !limiter.Allow("client-1") {
		t.Fatal("expected first client to be allowed")
	}

	if limiter.Allow("client-1") {
		t.Fatal("expected first client to be rate limited")
	}

	if !limiter.Allow("client-2") {
		t.Fatal("expected second client to have its own bucket")
	}
}

func TestLimiter_Allow_ConcurrentAccess(t *testing.T) {
	limiter := New(100, 100)

	const requests = 100

	var wg sync.WaitGroup
	wg.Add(requests)

	for range requests {
		go func() {
			defer wg.Done()

			_ = limiter.Allow("client-1")
		}()
	}

	wg.Wait()
}
