package errors

import (
	"fmt"
	"io"
)

type fundamental struct {
	msg string
	*stack
	cause error
}

func (f *fundamental) Msg() string {
	return f.msg
}

func (f *fundamental) Cause() error { return f.cause }

func (f *fundamental) Error() string {
	if f.cause == nil {
		return f.msg
	}

	return fmt.Sprintf("%s\nCaused by: %+v\n", f.msg, f.cause)
}

func (f *fundamental) Format(s fmt.State, verb rune) {
	io.WriteString(s, f.msg)
	f.stack.Format(s, verb)
	if f.cause != nil {
		fmt.Fprintf(s, "\nCaused by: %+v\n", f.cause)
	}
}

// Error returns an error with the supplied message.
// Error also records the stack trace at the point it was called.
func Error(message string) error {
	return &fundamental{
		msg:   message,
		stack: callers(),
	}
}

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
// Errorf also records the stack trace at the point it was called.
func Errorf(format string, args ...interface{}) error {
	return &fundamental{
		msg:   fmt.Sprintf(format, args...),
		stack: callers(),
	}
}

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	return &fundamental{
		msg:   message,
		cause: err,
		stack: callers(),
	}
}

// Wrapf returns an error annotating err with a stack trace
// at the point Wrapf is called, and the format specifier.
// If err is nil, Wrapf returns nil.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return &fundamental{
		msg:   fmt.Sprintf(format, args...),
		cause: err,
		stack: callers(),
	}
}
