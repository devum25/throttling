package main

import (
	"context"
	"time"
)

// Effector is the function which you want to throttle
type Effector func(context context.Context) (string,error)

// Throttled wraps and effector. It accepts the same parameters, plus a "UID" strings that represents a called ID
// It returns the same, plus a bool that's true if the call is not throttled
type Throttled func(context context.Context,uID string) (bool,string,error)

// A bucket tracks the request associated with a UID
type Bucket struct{
	tokens uint
	time time.Time
}

// Throttle returns an Effector function, and returns a Throttled function with a per-UID token bucket with a capacity of max
// that refills at the rate of refill token ever d.
func Throttle(e Effector,max uint,refill uint,d time.Duration) Throttled{
	buckets := map[string]*Bucket{}

	return func(context context.Context, uID string) (bool, string, error) {
        b := buckets[uID]

		if b == nil{
			buckets[uID] = &Bucket{tokens: max-1,time:time.Now()}

			str,err := e(context)
			return true,str,err
		}

		// calculate how many tokens we have based on time pass since the previous request.
		refillInterval := uint(time.Since(b.time)/d)
		tokensAdded := refill * refillInterval
		currentTokens := b.tokens + tokensAdded

		if currentTokens < 1{
			return false,"",nil
		}

 // If we've refilled our bucket, we can restart the clock.
 // Otherwise, we figure out when the most recent tokens were added.
 		if currentTokens > max {
			b.time = time.Now()
			b.tokens = max - 1
			} else {
				deltaTokens := currentTokens - b.tokens
				deltaRefills := deltaTokens / refill
				deltaTime := time.Duration(deltaRefills) * d
				b.time = b.time.Add(deltaTime)
				b.tokens = currentTokens - 1
	}
	str, err := e(context)
	return true, str, err
	}
}
