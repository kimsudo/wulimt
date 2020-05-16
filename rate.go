package wulimt

import (
	"context"
	"fmt"
	"time"

	"github.com/1set/gut/yrand"
	limiter "github.com/ulule/limiter/v3"
)

type RateLimiter struct {
	n string
	l *limiter.Limiter
}

var limiters = map[string]*RateLimiter{}

func (rl RateLimiter) String() string {
	return fmt.Sprintf("{limiter:%s, rate:%v/%v}", rl.n, rl.l.Rate.Limit, rl.l.Rate.Period)
}

// Wait blocks current goroutine to ensure the rate limit requirement is fulfilled.
func (rl *RateLimiter) Wait() {
	key := "default"
	for {
		ctx, _ := rl.l.Get(context.TODO(), key)
		if ctx.Reached {
			interval := time.Unix(ctx.Reset, 0).Sub(time.Now())
			if interval > 0 {
				timer := time.NewTimer(interval)
				<-timer.C
			} else {
				time.Sleep(200 * time.Millisecond)
			}
		} else {
			break
		}
	}
}

const (
	maxBufferRate = 0.4
)

// WaitMore blocks current goroutine to ensure the rate limit requirement is fulfilled, and wait a little longer.
func (rl *RateLimiter) WaitMore() {
	key := "default"
	for {
		ctx, _ := rl.l.Get(context.TODO(), key)
		if ctx.Reached {
			interval := time.Unix(ctx.Reset, 0).Sub(time.Now())
			if interval > 0 {
				bufferRate, _ := yrand.Float64()
				interval = time.Duration((bufferRate*maxBufferRate + 1) * float64(interval))
				timer := time.NewTimer(interval)
				<-timer.C
			} else {
				time.Sleep(400 * time.Millisecond)
			}
		} else {
			break
		}
	}
}

// Estimate returns average waiting time of the rate limiter.
func (rl *RateLimiter) Estimate() time.Duration {
	return time.Duration((float64(rl.l.Rate.Period) * (1 + 0.5*maxBufferRate)) / float64(rl.l.Rate.Limit))
}
