package err

type RuntimeError struct {
	*AbstractError
}

func NewRuntimeError(message string) *RuntimeError {
	return &RuntimeError{
		AbstractError: NewAbstractError(message, nil, nil),
	}
}

func NewRuntimeErrorWithCause(message string, cause any) *RuntimeError {
	return &RuntimeError{
		AbstractError: NewAbstractError(message, cause, nil),
	}
}

func NewRuntimeErrorWithStack(message string, cause any, stackTrace []uintptr) *RuntimeError {
	return &RuntimeError{
		AbstractError: NewAbstractError(message, cause, stackTrace),
	}
}
