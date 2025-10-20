package optional

import (
	"fmt"

	"github.com/go-external-config/go/lang"
)

type Optional[T any] struct {
	value   T
	present bool
	err     error
}

func OfNilable[T any](value T) *Optional[T] {
	return &Optional[T]{
		value:   value,
		present: !lang.IsNil(value)}
}

func OfEmpty[T any]() *Optional[T] {
	return &Optional[T]{present: false}
}

func OfValue[T any](value T) *Optional[T] {
	return &Optional[T]{
		value:   value,
		present: true}
}

func OfCommaOk[T any](value T, ok bool) *Optional[T] {
	if !ok {
		return OfEmpty[T]()
	}
	return OfValue(value)
}

func OfCommaErr[T any](value T, e error) *Optional[T] {
	if e != nil {
		return &Optional[T]{
			err: e,
		}
	}
	return OfValue(value)
}

func OfEntry[K comparable, V any](m map[K]V, key K) *Optional[V] {
	value, ok := m[key]
	if !ok {
		return OfEmpty[V]()
	}
	return OfValue(value)
}

func (o *Optional[T]) Present() bool {
	return o.present
}

func (o *Optional[T]) Value() T {
	o.panicIfEmpty("No value present")
	return o.value
}

func (o *Optional[T]) OrElse(value T) T {
	return lang.If(o.present, o.value, value)
}

func (o *Optional[T]) OrElseOptional(other *Optional[T]) *Optional[T] {
	return lang.If(o.present, o, other)
}

func (o *Optional[T]) OrElsePanic(format string, a ...any) T {
	o.panicIfEmpty(format, a...)
	return o.value
}

func (o *Optional[T]) panicIfEmpty(format string, a ...any) {
	if o.err != nil {
		panic(fmt.Errorf("%s\nCaused by: %v", fmt.Sprintf(format, a...), o.err.Error()))
	} else if !o.present {
		panic(fmt.Sprintf(format, a...))
	}
}
