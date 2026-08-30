package ratelimiting

import (
	"sync"
	"time"
)

type TokenBucket struct {
	Tokens         float64
	Capacity       int
	LastRefillTime time.Time
	RefillRate     float64
	mu             sync.RWMutex
}

func NewTokenBucket(capacity int, refillRate int) *TokenBucket {
	return &TokenBucket{
		Tokens:         float64(capacity),
		Capacity:       capacity,
		LastRefillTime: time.Now(),
		RefillRate:     float64(refillRate),
	}
}

func (t *TokenBucket) CheckConsumption() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.Tokens > 0 {
		t.Tokens--
		return true
	}
	return false
}

func (t *TokenBucket) Refill() {
	now := time.Now().UTC()
	elapsedTime := now.Sub(t.LastRefillTime).Seconds()
	tokensToAdd := elapsedTime * t.RefillRate
	if tokensToAdd > 0 {
		t.Tokens = min(tokensToAdd+t.Tokens, float64(t.Capacity))
		t.LastRefillTime = now
	}

}
