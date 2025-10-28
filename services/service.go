package services

import "errors"

var (
	ErrInvalidColor       = errors.New("invalid color")
	ErrInvalidDeviceSetup = errors.New("invalid device setup")
	ErrInvalidDevice      = errors.New("invalid device")
)

type Scannable interface {
	Scan(dest ...any) error
}
