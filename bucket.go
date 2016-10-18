package tokenbucket

import (
	"sync/atomic"
	"time"
)

// TokenBucket is an implementation of a token bucket.
type TokenBucket interface {
	// Consume attempts to consume the specified amount of tokens and
	// reports whether the operation was successful.
	// A change to the bucket is only made if Consume returns true.
	Consume(tokens uint64) bool
}

func New(rate, burstSize uint64) TokenBucket {
	return &tokenBucket{
		t:            0,
		timePerToken: 1000000 / rate,
		timePerBurst: burstSize * (1000000 / rate),
	}
}

type tokenBucket struct {
	t            uint64
	timePerToken uint64
	timePerBurst uint64
}

func (b *tokenBucket) Consume(tokens uint64) bool {
	timeNeeded := tokens * b.timePerToken
	oldTime := atomic.LoadUint64(&b.t)
	newTime := oldTime

	for {
		now := uint64(time.Now().UnixNano() / 1000)
		minTime := now - b.timePerBurst

		// Take into account burst size.
		// It is pretty unlikely that this is taken in any but the first
		// iteration of the for loop.
		if minTime > oldTime {
			newTime = minTime
		}

		// Now shift by the time needed.
		newTime += timeNeeded
		// Check if too many tokens.
		if newTime > now {
			return false
		}

		// CAS with the old value.
		if atomic.CompareAndSwapUint64(&b.t, oldTime, newTime) {
			return true
		}

		// Otherwise load old value and try again.
		oldTime = atomic.LoadUint64(&b.t)
		newTime = oldTime
	}

	return false
}
