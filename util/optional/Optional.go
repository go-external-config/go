package optional

import (
	"fmt"

	"github.com/go-errr/go/err"
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

func (this *Optional[T]) Present() bool {
	return this.present
}

func (this *Optional[T]) Value() T {
	this.panicIfEmpty("No value present")
	return this.value
}

func (this *Optional[T]) OrElse(value T) T {
	return lang.If(this.present, this.value, value)
}

func (this *Optional[T]) OrElseOptional(other *Optional[T]) *Optional[T] {
	return lang.If(this.present, this, other)
}

func (this *Optional[T]) OrElsePanic(format string, a ...any) T {
	this.panicIfEmpty(format, a...)
	return this.value
}

func (this *Optional[T]) panicIfEmpty(format string, a ...any) {
	if this.err != nil {
		panic(err.NewRuntimeExceptionFrom(fmt.Sprintf(format, a...), this.err))
	} else if !this.present {
		panic(err.NewNoSuchElementException(fmt.Sprintf(format, a...)))
	}
}
