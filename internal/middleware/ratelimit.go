package middleware

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

// Simple in-memory rate limiter for demonstration.
// In production, use Redis-based rate limiting for distributed systems.
type RateLimiter struct {
	ips map[string]*rate.Limiter
	mu  sync.Mutex
	r   rate.Limit
	b   int
}

func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		ips: make(map[string]*rate.Limiter),
		r:   r,
		b:   b,
	}
}

func (rl *RateLimiter) GetLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.r, rl.b)
		rl.ips[ip] = limiter
	}

	return limiter
}

func RateLimit(rl *RateLimiter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			limiter := rl.GetLimiter(ip)
			if !limiter.Allow() {
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"error": "Too many requests",
				})
			}
			return next(c)
		}
	}
}
