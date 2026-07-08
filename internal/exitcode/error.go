package exitcode

import "errors"

type CodedError struct {
	Err  error
	Code int
}

func (e *CodedError) Error() string {
	return e.Err.Error()
}

func (e *CodedError) Unwrap() error {
	return e.Err
}

func New(code int, err error) *CodedError {
	return &CodedError{Err: err, Code: code}
}

func Is(err error) bool {
	var coded *CodedError
	return errors.As(err, &coded)
}

func Code(err error) int {
	var coded *CodedError
	if errors.As(err, &coded) {
		return coded.Code
	}
	return ErrGeneral
}
