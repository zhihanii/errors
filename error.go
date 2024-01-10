package errors

import "fmt"

func WithCode(code int, format string, args ...any) error {
	return &codeError{
		err: fmt.Errorf(format, args...),
		code: code,
		stack: callers(),
	}
}

func Wrap(err error, code int, format string, args ...any) error {
	if err == nil {
		return nil
	}
	
	return &codeError{
		err: fmt.Errorf(format, args...),
		cause: err,
		code: code,
		stack: callers(),
	}
}

type codeError struct {
	err   error
	cause error
	code  int
	*stack
}

func (e *codeError) Error() string { return fmt.Sprintf("%v", e) }

func (e *codeError) Cause() error { return e.cause }

func (e *codeError) Unwrap() error { return e.cause }
