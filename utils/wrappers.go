package utils

import (
	"errors"
	"fmt"
)

// WrapError wraps an error with additional context, avoiding "failed to" patterns
func WrapError(err error, format string, a ...any) error {
	msg := fmt.Sprintf(format, a...)
	if err == nil {
		return errors.New(msg)
	}
	// Format the wrapped error with a concise message that starts with lowercase
	return fmt.Errorf("%s: %v", msg, err)
}
