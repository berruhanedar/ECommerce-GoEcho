package middleware

import (
	"fmt"
	"strings"
	"tesodev-korpes/pkg/customError"

	"github.com/labstack/echo/v4"
)

func AuthorizationMiddleware(allowedRoles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole, ok := c.Get("userRole").(string)
			fmt.Println(userRole)
			if !ok || userRole == "" {
				return customError.NewBadRequest(customError.EmptyRole)
			}

			for _, allowed := range allowedRoles {
				if strings.EqualFold(userRole, allowed) {
					return next(c)
				}
			}

			return customError.NewForbidden(customError.ForbiddenAccess)
		}
	}
}
