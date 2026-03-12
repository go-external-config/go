package err

import (
	"fmt"
)

type AbstractError struct {
	message string
	cause   error
	stack   []uintptr
}

func NewAbstractError(message string, cause any, stackTrace []uintptr) *AbstractError {
	var err error
	switch v := cause.(type) {
	case nil:
		err = nil
	case error:
		err = v
	default:
		err = fmt.Errorf("%v", v)
	}
	if stackTrace == nil {
		stackTrace = StackTrace(2)
	}

	return &AbstractError{
		message: message,
		cause:   err,
		stack:   stackTrace,
	}
}

func (this *AbstractError) Error() string {
	return this.message
}

func (this *AbstractError) Unwrap() error {
	return this.cause
}

func (this *AbstractError) StackTrace() []uintptr {
	return this.stack
}
