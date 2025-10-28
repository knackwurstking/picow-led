package services

import "errors"

var (
	ErrInvalidColor       = errors.New("invalid color")
	ErrInvalidDeviceSetup = errors.New("invalid device setup")
	ErrInvalidDevice      = errors.New("invalid device")
	ErrInvalidDeviceID    = errors.New("invalid device ID")
	ErrNotFound           = errors.New("not found")
)

type Scannable interface {
	Scan(dest ...any) error
}
