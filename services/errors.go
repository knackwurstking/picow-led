package services

import (
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrInvalidColor       = errors.New("invalid color")
	ErrInvalidDeviceSetup = errors.New("invalid device setup")
	ErrInvalidDevice      = errors.New("invalid device")
	ErrInvalidGroup       = errors.New("invalid group")

	ErrInvalidGroupDeviceID = errors.New("invalid device ID")
)

type NotFoundError struct {
	Message string
}

func NewNotFoundError(format string, a ...any) error {
	return &NotFoundError{Message: fmt.Sprintf(format, a...)}
}

func (e NotFoundError) Error() string {
	return e.Message
}

func IsNotFoundError(err error) bool {
	return errors.Is(err, &NotFoundError{})
}

func HandleSqlError(err error) error {
	if err == sql.ErrNoRows {
		return NewNotFoundError("not found: %v", err)
	}

	return err
}
