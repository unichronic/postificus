package middleware

import (
	"net/http"
	"time"

	"postificus/internal/storage"

	"github.com/labstack/echo/v4"
)

func RateLimit(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if storage.RedisClient == nil {
			return next(c)
		}
		key := "rl:" + c.RealIP()
		ctx := c.Request().Context()
		count, _ := storage.RedisClient.Incr(ctx, key).Result()
		if count == 1 {
			storage.RedisClient.Expire(ctx, key, time.Minute)
		}
		if count > 60 {
			return c.JSON(http.StatusTooManyRequests, map[string]string{"error": "Too many requests"})
		}
		return next(c)
	}
}
