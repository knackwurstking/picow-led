package handlers

import (
	"errors"
	"fmt"
)

type ValidationError struct {
	Message string
}

func NewValidationError(message string, args ...any) error {
	return &ValidationError{
		Message: fmt.Sprintf(message, args...),
	}
}

func (e ValidationError) Error() string {
	return e.Message
}

type DatabaseError struct {
	Message string
}

func NewDatabaseError(message string, args ...any) error {
	return &DatabaseError{
		Message: fmt.Sprintf(message, args...),
	}
}

func (e DatabaseError) Error() string {
	return e.Message
}

func IsValidationError(err error) bool {
	return errors.Is(err, &ValidationError{})
}

func IsDatabaseError(err error) bool {
	return errors.Is(err, &DatabaseError{})
}
