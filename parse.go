package wulimt

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	limiter "github.com/ulule/limiter/v3"
)

func ParseRate(formatted string) (*limiter.Rate, error) {
	rate := &limiter.Rate{}

	values := strings.Split(formatted, "-")
	if len(values) != 2 {
		return nil, fmt.Errorf("incorrect format '%s'", formatted)
	}

	periodMap := map[string]time.Duration{
		"S": time.Second,    // Second
		"M": time.Minute,    // Minute
		"H": time.Hour,      // Hour
		"D": time.Hour * 24, // Day
	}

	limit, period := strings.TrimSpace(values[0]), strings.TrimSpace(values[1])

	numPeriod := time.Duration(0)
	if lp := len(period); lp == 0 {
		return nil, fmt.Errorf("incorrect period '%s': blank", period)
	} else {
		base := int64(1)
		if lp > 1 {
			pn, err := strconv.ParseInt(period[0:lp-1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("incorrect period '%s': %v", limit, err)
			}
			base = pn
		}

		periodUnit := strings.ToUpper(string(period[lp-1]))
		duration, ok := periodMap[periodUnit]
		if !ok {
			return nil, fmt.Errorf("incorrect period unit '%s'", period)
		}

		numPeriod = time.Duration(base) * duration
	}

	numLimit, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("incorrect limit '%s': %v", limit, err)
	} else if numLimit < 1 {
		return nil, fmt.Errorf("invalid limit '%s': %d", limit, numLimit)
	}

	rate = &limiter.Rate{
		Formatted: formatted,
		Period:    numPeriod,
		Limit:     numLimit,
	}
	return rate, nil
}
