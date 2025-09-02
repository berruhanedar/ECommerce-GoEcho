package middleware

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type contextKey string

const correlationIDKey contextKey = "CorrelationID"

func CorrelationIdMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			correlationID := req.Header.Get("X-Correlation-ID")
			if correlationID == "" {
				correlationID = uuid.New().String()
			}

			res.Header().Set("X-Correlation-ID", correlationID)

			ctx := context.WithValue(req.Context(), correlationIDKey, correlationID)
			c.SetRequest(req.WithContext(ctx))
			c.Set(string(correlationIDKey), correlationID)

			return next(c)
		}
	}
}
