package ratelimit

import (
	"sync"
	"time"
)

//a bucket has tokens (maybe keep the size dynamic?)
// refill the tokens
// set a timing window for refill

type tokenBucket struct {
	currentTokens int
	lastUpdatedAT time.Time
}

type RateLimiterConfig struct {
	visitor          map[any]*tokenBucket
	maxAllowedTokens int
	refillRate       float64
	mu               sync.Mutex
}

//something to set the config for what the request save type is

// validating each request
func Validate(rq *RateLimiterConfig, ip string) bool {
	rq.mu.Lock()
	defer rq.mu.Unlock()

	bucket, exists := rq.visitor[ip]
	//new user if it does not exist so we can return without any other calculations
	if !exists {
		tb := &tokenBucket{
			currentTokens: rq.maxAllowedTokens - 1,
			lastUpdatedAT: time.Now(),
		}

		rq.visitor[ip] = tb

		return true
	}

	elapsed := time.Since(bucket.lastUpdatedAT).Seconds()
	tokensToAdd := int(elapsed * rq.refillRate)

	if tokensToAdd > 0 {
		bucket.currentTokens += tokensToAdd
		if bucket.currentTokens > rq.maxAllowedTokens {
			bucket.currentTokens = rq.maxAllowedTokens
		}
		bucket.lastUpdatedAT = time.Now()
	}

	//checking if currentTokens are already not zero before subtracting
	if bucket.currentTokens > 0 {
		bucket.currentTokens -= 1
		return true
	}

	return false
}

// constructor func
func NewRateLimterConfig(maxAllowedTokens int, refillRate float64) *RateLimiterConfig {
	return &RateLimiterConfig{
		visitor:          map[any]*tokenBucket{},
		maxAllowedTokens: maxAllowedTokens,
		refillRate:       refillRate,
		mu:               sync.Mutex{},
	}
}
