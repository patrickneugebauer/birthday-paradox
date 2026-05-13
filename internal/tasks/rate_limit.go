package tasks

import (
	"net/http"
	"strconv"
	"time"
)

type RateLimit struct {
	Limit     int
	Remaining int
	Reset     time.Time
}

func parseRateLimitHeaders(h http.Header) RateLimit {
	rl := RateLimit{
		Limit:     60,
		Remaining: 60,
		Reset:     time.Now().Add(time.Hour),
	}

	if limit := h.Get("X-RateLimit-Limit"); limit != "" {
		if v, err := strconv.Atoi(limit); err == nil {
			rl.Limit = v
		}
	}

	if remaining := h.Get("X-RateLimit-Remaining"); remaining != "" {
		if v, err := strconv.Atoi(remaining); err == nil {
			rl.Remaining = v
		}
	}

	if reset := h.Get("X-RateLimit-Reset"); reset != "" {
		if v, err := strconv.ParseInt(reset, 10, 64); err == nil {
			rl.Reset = time.Unix(v, 0)
		}
	}

	return rl
}
