package services

import "errors"

var (
	ErrorNotFound   = errors.New("not found")
	ErrorValidation = errors.New("validation failed")
)
