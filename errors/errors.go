package errors

import (
	"errors"
	"fmt"
)

func Wrap(err error, format string, a ...any) error {
	msg := fmt.Sprintf(format, a...)
	if err == nil {
		return errors.New(msg)
	}
	return fmt.Errorf("%s: %v", msg, err)
}
