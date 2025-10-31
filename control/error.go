package control

import "errors"

var (
	ErrNotConnected = errors.New("not connected")
	ErrNoData       = errors.New("no data")
)
