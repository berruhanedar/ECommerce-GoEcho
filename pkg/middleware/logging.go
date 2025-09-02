package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		start := time.Now()

		req := c.Request()
		res := c.Response()
		method := req.Method
		path := req.URL.Path
		ip := c.RealIP()
		userAgent := req.UserAgent()
		correlationID := c.Response().Header().Get("X-Correlation-ID")

		logrus.WithFields(logrus.Fields{
			"event":         "request",
			"timestamp":     start.Format(time.RFC3339),
			"method":        method,
			"path":          path,
			"ip":            ip,
			"user_agent":    userAgent,
			"correlationID": correlationID,
		}).Info("Incoming request")

		err := next(c)

		stop := time.Now()
		latency := stop.Sub(start)
		status := res.Status

		responseLog := logrus.WithFields(logrus.Fields{
			"event":         "response",
			"timestamp":     stop.Format(time.RFC3339),
			"method":        method,
			"path":          path,
			"status":        status,
			"latency":       latency.String(),
			"ip":            ip,
			"correlationID": correlationID,
		})

		if err != nil {
			responseLog.WithField("error", err.Error()).Error("Request failed")
		} else {
			responseLog.Info("Request completed")
		}

		return err

	}
}
