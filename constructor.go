package wulimt

import (
	"fmt"
	"time"

	limiter "github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func Get(name string) *RateLimiter {
	if rl, ok := limiters[name]; ok {
		return rl
	} else {
		return nil
	}
}

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
		}
		limiters[name] = rl
		return rl
	}
}

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
		}
		limiters[name] = rl
		return rl
	}
}