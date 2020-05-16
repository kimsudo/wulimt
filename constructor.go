// Package wulimt is a wrapper for ulule/limiter to provide simple rate limiter for client-side app.
package wulimt

import (
	"fmt"
	"time"

	limiter "github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// Get returns rate limiter by name, nil will be returned if not exists.
func Get(name string) *RateLimiter {
	if rl, ok := limiters[name]; ok {
		return rl
	} else {
		return nil
	}
}

// Get returns rate limiter by name, or create an new one if not exists.
func GetOrNew(name string, periodSec, times int64) *RateLimiter {
	if rl, ok := limiters[name]; ok {
		return rl
	} else {
		rl = &RateLimiter{
			n: name,
			l: limiter.New(memory.NewStore(), limiter.Rate{
				Period: time.Duration(periodSec) * time.Second,
				Limit:  times,
			}),
			b: defaultBufferRate,
		}
		limiters[name] = rl
		return rl
	}
}

// Get returns rate limiter by name, or create an new one if not exists. It panics if fails to parse formatted rate string.
func GetOrParseNew(name, expr string) *RateLimiter {
	if rl, ok := limiters[name]; ok {
		return rl
	} else {
		rate, err := ParseRate(expr)
		if err != nil {
			panic(err)
		} else if rate == nil {
			panic(fmt.Errorf("got nil limiter.Rate"))
		}

		rl = &RateLimiter{
			n: name,
			l: limiter.New(memory.NewStore(), *rate),
			b: defaultBufferRate,
		}
		limiters[name] = rl
		return rl
	}
}
