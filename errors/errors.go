package errors

import (
	"errors"
	"fmt"
)

// AppError represents a generic application error with context
type AppError struct {
	Err     error
	Code    string
	Message string
	Context map[string]any
}

// Error implements the error interface for AppError
func (e *AppError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Err.Error()
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// Is checks if the error is of the specified type
func (e *AppError) Is(target error) bool {
	return errors.Is(e.Err, target)
}

// New creates a new AppError with the provided error and context
func New(code, message string, err error, context map[string]any) *AppError {
	return &AppError{
		Err:     err,
		Code:    code,
		Message: message,
		Context: context,
	}
}

// NewWithCode creates a new AppError with just the code and message
func NewWithCode(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Error codes for the application
const (
	CodeInvalidFlags        = "INVALID_FLAGS"
	CodeInvalidLogFormat    = "INVALID_LOG_FORMAT"
	CodeDatabaseConnection  = "DATABASE_CONNECTION_ERROR"
	CodeDatabasePing        = "DATABASE_PING_ERROR"
	CodeInvalidDatabasePath = "INVALID_DATABASE_PATH"
	CodeDatabaseTables      = "DATABASE_TABLES_ERROR"
	CodeSetupDevices        = "SETUP_DEVICES_ERROR"
	CodeServerStart         = "SERVER_START_ERROR"
	CodeDeviceNotFound      = "DEVICE_NOT_FOUND"
	CodeGroupNotFound       = "GROUP_NOT_FOUND"
	CodeColorInvalid        = "COLOR_INVALID"
	CodeDeviceSetupInvalid  = "DEVICE_SETUP_INVALID"
	CodeInvalidAddress      = "INVALID_ADDRESS"
	CodeInvalidPort         = "INVALID_PORT"
	CodeDeviceControlError  = "DEVICE_CONTROL_ERROR"
	CodeColorNotFound       = "COLOR_NOT_FOUND"
	CodeGroupDeviceInvalid  = "GROUP_DEVICE_INVALID"
)

// Specific error types
var (
	ErrInvalidFlags        = NewWithCode(CodeInvalidFlags, "Invalid command flags")
	ErrInvalidLogFormat    = NewWithCode(CodeInvalidLogFormat, "Invalid log format")
	ErrDatabaseConnection  = NewWithCode(CodeDatabaseConnection, "Failed to connect to database")
	ErrDatabasePing        = NewWithCode(CodeDatabasePing, "Failed to ping database")
	ErrInvalidDatabasePath = NewWithCode(CodeInvalidDatabasePath, "Database path is required")
	ErrDatabaseTables      = NewWithCode(CodeDatabaseTables, "Failed to create database tables")
	ErrSetupDevices        = NewWithCode(CodeSetupDevices, "Failed to set up devices")
	ErrServerStart         = NewWithCode(CodeServerStart, "Failed to start server")
	ErrDeviceNotFound      = NewWithCode(CodeDeviceNotFound, "Device not found")
	ErrGroupNotFound       = NewWithCode(CodeGroupNotFound, "Group not found")
	ErrColorInvalid        = NewWithCode(CodeColorInvalid, "Invalid color format")
	ErrDeviceSetupInvalid  = NewWithCode(CodeDeviceSetupInvalid, "Invalid device setup")
	ErrInvalidAddress      = NewWithCode(CodeInvalidAddress, "Invalid device address")
	ErrInvalidPort         = NewWithCode(CodeInvalidPort, "Invalid device port")
	ErrDeviceControlError  = NewWithCode(CodeDeviceControlError, "Device control error")
	ErrColorNotFound       = NewWithCode(CodeColorNotFound, "Color not found in database")
	ErrGroupDeviceInvalid  = NewWithCode(CodeGroupDeviceInvalid, "Invalid device in group")
)

// IsAppError checks if an error is of type AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetErrorCode returns the error code if the error is an AppError
func GetErrorCode(err error) string {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code
	}
	return ""
}

// GetErrorContext returns the error context if the error is an AppError
func GetErrorContext(err error) map[string]any {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Context
	}
	return nil
}

// Wrap wraps an error with additional context and returns an AppError
func Wrap(err error, code, message string, context map[string]any) *AppError {
	if err == nil {
		return nil
	}

	// If it's already an AppError, just add the context and return it
	if appErr, ok := err.(*AppError); ok {
		// Merge contexts if both have contexts
		if appErr.Context != nil && context != nil {
			merged := make(map[string]any)
			for k, v := range appErr.Context {
				merged[k] = v
			}
			for k, v := range context {
				merged[k] = v
			}
			appErr.Context = merged
		} else if context != nil {
			appErr.Context = context
		}

		// If the error code is the same or we're not overriding, return as-is
		if appErr.Code == "" {
			appErr.Code = code
		}

		return appErr
	}

	return &AppError{
		Err:     err,
		Code:    code,
		Message: message,
		Context: context,
	}
}

// WrapWithCode wraps an error with a specific code and message
func WrapWithCode(err error, code, message string) *AppError {
	if err == nil {
		return nil
	}

	return &AppError{
		Err:     err,
		Code:    code,
		Message: message,
		Context: nil,
	}
}

// FormatError formats an error with additional context for logging
func FormatError(err error, context map[string]any) string {
	if err == nil {
		return "no error"
	}

	if appErr, ok := err.(*AppError); ok {
		if appErr.Message != "" {
			return fmt.Sprintf("%s: %v", appErr.Message, appErr.Err)
		}
		return appErr.Error()
	}

	return err.Error()
}
