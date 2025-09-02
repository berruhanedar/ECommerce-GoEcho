package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/labstack/echo/v4"
)

func RecoveryMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if r := recover(); r != nil {

				c.Logger().Errorf("PANIC recovered: %v\n%s", r, debug.Stack())

				_ = c.JSON(http.StatusInternalServerError, echo.Map{
					"error":   "Internal server error",
					"message": "A system error occurred. Please try again later.",
				})
			}
		}()
		return next(c)
	}
}
