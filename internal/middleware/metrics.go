package middleware

import (
	"strconv"
	"time"

	"postificus/internal/metrics"

	"github.com/labstack/echo/v4"
)

func PrometheusMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		err := next(c)

		status := c.Response().Status
		duration := time.Since(start).Seconds()

		path := c.Path()
		method := c.Request().Method

		// Normalize path if needed (e.g., replace IDs with :id)
		// Echo paths are already parameterized (e.g., /users/:id), so we use c.Path()

		metrics.HTTPRequestsTotal.WithLabelValues(method, path, strconv.Itoa(status)).Inc()
		metrics.HTTPRequestDuration.WithLabelValues(method, path).Observe(duration)

		return err
	}
}
