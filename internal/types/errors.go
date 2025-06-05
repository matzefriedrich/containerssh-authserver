package types

import "errors"

type AuthServerError struct {
	cause error
	Msg   string
}

func (f AuthServerError) Error() string {
	return f.Msg
}

func (f AuthServerError) Unwrap() error {
	return f.cause
}

func (f AuthServerError) Is(err error) bool {
	return f.Error() == err.Error()
}

type AuthServerAggregateError struct {
	errors []error
	Msg    string
}

type AuthServerErrorFunc func(e error)

func (f AuthServerAggregateError) Errors() []error {
	return f.errors
}

func (f AuthServerAggregateError) Error() string {
	return f.Msg
}

func (f AuthServerAggregateError) Is(err error) bool {
	if f.Error() == err.Error() {
		return true
	}
	for _, cause := range f.errors {
		if errors.Is(err, cause) {
			return true
		}
	}
	return false
}

func WithAggregatedCause(errs ...error) AuthServerErrorFunc {
	return func(e error) {
		var funqErr *AuthServerError
		ok := errors.As(e, &funqErr)
		if ok {
			funqErr.cause = &AuthServerAggregateError{
				errors: errs,
				Msg:    "one or more errors occurred",
			}
		}
	}
}
