package lang

import (
	"fmt"
	"reflect"
)

func If[T any](cond bool, v1, v2 T) T {
	if cond {
		return v1
	}
	return v2
}

func IsNil(value any) bool {
	if value == nil {
		return true
	}
	defer func() { recover() }()
	return reflect.ValueOf(value).IsNil()
}

func FirstNonEmpty[T comparable](values ...T) T {
	empty := *new(T)
	for _, value := range values {
		if value != empty {
			return value
		}
	}
	return empty
}

func AssertState(expression bool, format string, args ...any) {
	if !expression {
		panic(fmt.Sprintf(format, args...))
	}
}
