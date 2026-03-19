package errors

import (
	"errors"
	"fmt"
)

// AppError is the standard application error carrying a code, human-readable message,
// optional wrapped cause, and optional stack trace.
type AppError struct {
	Code    Code
	Message string
	Err     error
	frames  []Frame
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// Stack returns captured stack frames, if any.
func (e *AppError) Stack() []Frame {
	return e.frames
}

// WithStack captures the current call stack and attaches it to the error.
func (e *AppError) WithStack() *AppError {
	e.frames = callers(2)
	return e
}

// New creates an AppError with the given code and message.
func New(code Code, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// Wrap creates an AppError that wraps an existing error.
func Wrap(code Code, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

// HasCode reports whether any error in err's tree has the given code.
func HasCode(err error, code Code) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == code
	}
	return false
}

// GetCode returns the Code of the first AppError found in err's tree.
// Returns CodeUnknown if err is nil or no AppError is found.
func GetCode(err error) Code {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}
	return CodeUnknown
}

// As is errors.As exposed for use without importing the standard errors package.
var As = errors.As

// Is is errors.Is exposed for use without importing the standard errors package.
var Is = errors.Is
