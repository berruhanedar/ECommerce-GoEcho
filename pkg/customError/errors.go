package customError

import (
	"fmt"

	"github.com/labstack/gommon/log"
)

type AppError struct {
	HTTPStatus int    `json:"-"`
	Code       int    `json:"code"`
	Message    string `json:"message"`
	Err        error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code: %v, message: %v, underlying_error: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("code: %v, message: %v", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func newAppError(key ErrorKey, err error, args ...interface{}) *AppError {
	details, ok := ErrorDefinitions[key]
	if !ok {
		log.Errorf("FATAL: Undefined error key requested: %s", key)
		details = ErrorDefinitions[InternalServerError]
	}

	msg := details.Message
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	return &AppError{
		HTTPStatus: details.StatusCode,
		Code:       details.TypeCode,
		Message:    msg,
		Err:        err,
	}
}

func NewNotFound(key ErrorKey) *AppError {
	return newAppError(key, nil)
}

func NewBadRequest(key ErrorKey) *AppError {
	return newAppError(key, nil)
}

func NewUnauthorized(key ErrorKey) *AppError {
	return newAppError(key, nil)
}

func NewForbidden(key ErrorKey) *AppError {
	return newAppError(key, nil)
}

func NewConflict(key ErrorKey, args ...interface{}) *AppError {
	return newAppError(key, nil, args...)
}

func NewUnprocessableEntity(key ErrorKey, err error) *AppError {
	return newAppError(key, err)
}

func NewInternal(key ErrorKey, err error) *AppError {
	log.Errorf("Internal error occurred with key %s: %v", key, err)
	return newAppError(key, err)
}

func NewValidation(key ErrorKey) *AppError {
	return newAppError(key, nil)
}

func LogErrorWithCorrelation(err error, correlationID string) {

	log.Printf("[CorrelationID: %s] %s", correlationID, err.Error())

}
func LogInfoWithCorrelation(message string, correlationID string) {
	log.Printf("[CorrelationID: %s] %s", correlationID, message)
}
