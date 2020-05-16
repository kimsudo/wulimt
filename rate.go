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
	b float64
}

const (
	defaultBufferRate = 0.4
)

var limiters = map[string]*RateLimiter{}

func (rl RateLimiter) String() string {
	return fmt.Sprintf("{limiter:%s, rate:%v/%v}", rl.n, rl.l.Rate.Limit, rl.l.Rate.Period)
}

// GetBufferRate returns buffer rate of the rate limiter.
func (rl *RateLimiter) GetBufferRate() float64 {
	return rl.b
}

// GetBufferRate sets buffer rate of the rate limiter, value <= 0 will disable time buffer in WaitMore() as well.
func (rl *RateLimiter) SetBufferRate(br float64) {
	rl.b = br
}

// Wait blocks current goroutine to ensure the rate limit requirement is fulfilled.
func (rl *RateLimiter) Wait() {
	waitOnLimiter(rl.l, -1)
}

// WaitMore blocks current goroutine to ensure the rate limit requirement is fulfilled, and wait a little longer.
func (rl *RateLimiter) WaitMore() {
	waitOnLimiter(rl.l, rl.b)
}

// Estimate returns average waiting time of the rate limiter.
func (rl *RateLimiter) Estimate() time.Duration {
	return time.Duration((float64(rl.l.Rate.Period) * (1 + 0.5*rl.b)) / float64(rl.l.Rate.Limit))
}

func waitOnLimiter(rl *limiter.Limiter, bufferRate float64) {
	key := "default"
	for {
		ctx, _ := rl.Get(context.TODO(), key)
		if ctx.Reached {
			interval := time.Unix(ctx.Reset, 0).Sub(time.Now())
			if interval > 0 {
				if bufferRate > 0 {
					rand, _ := yrand.Float64()
					interval = time.Duration((bufferRate*rand + 1) * float64(interval))
				}
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
