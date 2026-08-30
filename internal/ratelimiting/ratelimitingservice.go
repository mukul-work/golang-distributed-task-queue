package ratelimiting

import "sync"

type RateLimiter struct {
	Capacity   int
	RefillRate int
	Buckets    map[string]*TokenBucket
	mu         *sync.Mutex
}

func NewRateLimiterService(capacity, refillRate int) *RateLimiter {
	return &RateLimiter{
		Buckets:    make(map[string]*TokenBucket),
		Capacity:   capacity,
		RefillRate: refillRate,
		mu:         &sync.Mutex{},
	}
}
func (r *RateLimiter) Allow(keyID string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	bucket, exists := r.Buckets[keyID]
	if !exists {
		bucket = NewTokenBucket(r.Capacity, r.RefillRate)
		r.Buckets[keyID] = bucket
	}
	return bucket.CheckConsumption()
}
