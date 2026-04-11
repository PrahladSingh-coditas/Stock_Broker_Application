package models

import (
	"sync"
	"time"
)

var bucket *TokenBucket
var once sync.Once

type TokenBucket struct {
	Capacity     int64
	Tokens       int64
	Rate         int64
	LastRefilled time.Time
	Mut          sync.Mutex
}

func GetNewTokenBucket(capacity int64, tokensPerSecond int64) *TokenBucket {

	bucket = &TokenBucket{
		Capacity:     capacity,
		Tokens:       capacity,
		Rate:         tokensPerSecond,
		LastRefilled: time.Now(),
	}

	return bucket
}

func (tb *TokenBucket) refill() {
	currTime := time.Now()
	passedTime := currTime.Sub(tb.LastRefilled).Seconds()

	tb.Tokens += int64(passedTime) * tb.Rate

	if tb.Tokens > tb.Capacity {
		tb.Tokens = tb.Capacity
	}

	tb.LastRefilled = time.Now()
}

func (tb *TokenBucket) Allow() bool {
	tb.Mut.Lock()
	defer tb.Mut.Unlock()

	tb.refill()

	if tb.Tokens >= 1 {
		tb.Tokens--
		return true
	}
	return false
}
