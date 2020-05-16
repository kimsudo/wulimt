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

func (rl *RateLimiter) String() string {
	return fmt.Sprintf("{limiter:%s, rate:%v/%v}", rl.n, rl.l.Rate.Limit, rl.l.Rate.Period)
}

func (rl *RateLimiter) Wait() {
	key := "default"
	for {
		ctx, _ := rl.l.Get(context.TODO(), key)
		if ctx.Reached {
			interval := time.Unix(ctx.Reset, 0).Sub(time.Now())
			if interval > 0 {
				//log.Debug("wait for next period", zap.Duration("left", interval))
				timer := time.NewTimer(interval)
				<-timer.C
				//log.Debug("time is up, let's retry")
			} else {
				//log.Debug("not ready yet", zap.Duration("left", interval))
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

func (rl *RateLimiter) WaitMore() {
	key := "default"
	for {
		ctx, _ := rl.l.Get(context.TODO(), key)
		if ctx.Reached {
			interval := time.Unix(ctx.Reset, 0).Sub(time.Now())
			if interval > 0 {
				bufferRate, _ := yrand.Float64()
				fmt.Printf("~~~ original: %v, rate: %f, ", interval, bufferRate)
				interval = time.Duration((bufferRate*maxBufferRate + 1) * float64(interval))
				fmt.Printf("adjusted: %v\n", interval)
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

func (rl *RateLimiter) Estimate() time.Duration {
	return time.Duration((float64(rl.l.Rate.Period) * (1 + 0.5*maxBufferRate)) / float64(rl.l.Rate.Limit))
}
