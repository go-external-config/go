package lang

import "fmt"

type Optional[T any] struct {
	value   T
	present bool
	err     error
}

func OptionalEmpty[T any]() *Optional[T] {
	return &Optional[T]{present: false}
}

func OptionalValue[T any](value T) *Optional[T] {
	return &Optional[T]{
		value:   value,
		present: true}
}

func OptionalOfCommaOk[T any](value T, ok bool) *Optional[T] {
	if !ok {
		return OptionalEmpty[T]()
	}
	return OptionalValue(value)
}

func OptionalOfCommaErr[T any](value T, e error) *Optional[T] {
	if e != nil {
		return &Optional[T]{
			err: e,
		}
	}
	return OptionalValue(value)
}

func OptionalEntry[K comparable, V any](m map[K]V, key K) *Optional[V] {
	value, ok := m[key]
	if !ok {
		return OptionalEmpty[V]()
	}
	return OptionalValue(value)
}

func (o *Optional[T]) Present() bool {
	return o.present
}

func (o *Optional[T]) Value() T {
	return o.value
}

func (o *Optional[T]) OrElse(value T) T {
	return If(o.present, o.value, value)
}

func (o *Optional[T]) OrElsePanic(msg string) T {
	if o.err != nil {
		panic(fmt.Errorf("%s\nCaused by: %v", msg, o.err.Error()))
	} else if !o.present {
		panic(msg)
	}
	return o.value
}
