package app

import "errors"

type ErrorType string

const (
	ErrNotFound     ErrorType = "not_found"
	ErrConflict     ErrorType = "conflict"
	ErrInvalidInput ErrorType = "invalid_input"
	ErrInternal     ErrorType = "internal"
	ErrUnauthorized ErrorType = "unauthorized"
)

type AppError struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(errType ErrorType, message string, err error) *AppError {
	return &AppError{
		Type:    errType,
		Message: message,
		Err:     err,
	}
}

func Wrap(err error, errType ErrorType, message string) *AppError {
	return &AppError{
		Type:    errType,
		Message: message,
		Err:     err,
	}
}

func (e *AppError) Is(target error) bool {
	var appErr *AppError
	if errors.As(target, &appErr) {
		return appErr.Type == e.Type
	}
	return false
}
