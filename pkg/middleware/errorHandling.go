package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"tesodev-korpes/pkg/customError"
	"time"

	"github.com/labstack/echo/v4"
)

type APIErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	CorrelationID string `json:"correlationID"`
	Timestamp     string `json:"timestamp"`
}

func ErrorHandler() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}

			var httpErr *echo.HTTPError
			if errors.As(err, &httpErr) {

				fmt.Printf("[HTTP %d] %v\n", httpErr.Code, httpErr.Message)

				switch httpErr.Code {
				case http.StatusNotFound:
					return customError.NewNotFound(customError.UnknownFotFound)
				case http.StatusBadRequest:
					return customError.NewBadRequest(customError.UnknownBadRequest)
				case http.StatusInternalServerError:
					return customError.NewInternal(customError.UnknownServiceError, err)
				default:
					return customError.NewInternal(customError.FrameworkError, httpErr)
				}
			}

			appErr := toAppError(err)

			if c.Response().Committed {
				return nil
			}

			return c.JSON(appErr.HTTPStatus, buildAPIResponse(appErr, c))
		}
	}
}

func toAppError(err error) *customError.AppError {
	var appErr *customError.AppError

	if errors.As(err, &appErr) {
		return appErr
	}

	return customError.NewInternal(customError.InternalServerError, err)
}

func buildAPIResponse(err *customError.AppError, c echo.Context) APIErrorResponse {
	correlationID := c.Response().Header().Get(echo.HeaderXCorrelationID)
	if correlationID == "" {
		correlationID = "not-available"
	}

	var resp APIErrorResponse
	resp.Error.Code = err.Code
	resp.Error.Message = err.Message
	resp.CorrelationID = correlationID
	resp.Timestamp = time.Now().UTC().Format(time.RFC3339)
	return resp
}
