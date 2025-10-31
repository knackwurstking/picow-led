package services

import (
	"database/sql"
	"errors"
)

var (
	ErrInvalidColor       = errors.New("invalid color")
	ErrInvalidDeviceSetup = errors.New("invalid device setup")
	ErrInvalidDevice      = errors.New("invalid device")
	ErrInvalidDeviceID    = errors.New("invalid device ID")
	ErrNotFound           = errors.New("not found")
)

func HandleSqlError(err error) error {
	if err == sql.ErrNoRows {
		return ErrNotFound
	}

	return err
}
