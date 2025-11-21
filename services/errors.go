package services

import (
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrInvalidColor         = errors.New("invalid color provided")
	ErrInvalidDeviceSetup   = errors.New("invalid device setup provided")
	ErrInvalidDevice        = errors.New("invalid device provided")
	ErrInvalidGroup         = errors.New("invalid group provided")
	ErrInvalidGroupDeviceID = errors.New("invalid device ID in group")
)

// NotFoundError represents an error when a resource is not found
type NotFoundError struct {
	Message string
}

// NewNotFoundError creates a new NotFoundError with formatted message
func NewNotFoundError(msg string) error {
	return &NotFoundError{Message: msg}
}

// Error implements the error interface for NotFoundError
func (e NotFoundError) Error() string {
	return e.Message
}

// IsNotFoundError checks if an error is a NotFoundError
func IsNotFoundError(err error) bool {
	if e, ok := err.(*ServiceError); ok {
		err = e.Unwrap()
	}

	_, ok := err.(*NotFoundError)
	return ok
}

// HandleSqlError converts SQL errors to more descriptive service errors
func HandleSqlError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return NewServiceError("not found", err)
	default:
		return NewServiceError("database operation failed", err)
	}
}

// ServiceError wraps service errors with additional context
type ServiceError struct {
	Operation string
	Err       error
}

// NewServiceError creates a new ServiceError with operation context
func NewServiceError(operation string, err error) error {
	if err == nil {
		return nil
	}

	return &ServiceError{
		Operation: operation,
		Err:       err,
	}
}

// Error implements the error interface for ServiceError
func (e *ServiceError) Error() string {
	return fmt.Sprintf("service operation '%s': %v", e.Operation, e.Err)
}

// Unwrap returns the wrapped error
func (e *ServiceError) Unwrap() error {
	return e.Err
}
