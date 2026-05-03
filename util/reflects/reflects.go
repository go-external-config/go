package reflects

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/go-errr/go/err"
)

type Field struct {
	Owner    reflect.Type
	Field    reflect.StructField
	Index    []int
	Type     reflect.Type
	Value    reflect.Value
	TagName  string
	TagValue string
}

func ForEachTaggedField(target any, tagName string, fn func(Field)) {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Pointer || targetValue.IsNil() {
		panic(err.NewIllegalArgumentException(fmt.Sprintf("Target must be non-nil pointer to struct, got %T", target)))
	}

	structValue := targetValue.Elem()
	if structValue.Kind() != reflect.Struct {
		panic(err.NewIllegalArgumentException(fmt.Sprintf("Target must point to struct, got %s", structValue.Kind())))
	}

	structType := structValue.Type()
	for i := 0; i < structType.NumField(); i++ {
		structField := structType.Field(i)
		tagValue, ok := structField.Tag.Lookup(tagName)
		if !ok {
			continue
		}
		fn(Field{
			Owner:    structType,
			Field:    structField,
			Index:    structField.Index,
			Type:     structField.Type,
			Value:    Settable(structValue.Field(i)),
			TagName:  tagName,
			TagValue: tagValue,
		})
	}
}

func Settable(v reflect.Value) reflect.Value {
	if v.CanSet() {
		return v
	}
	ptr := unsafe.Pointer(v.UnsafeAddr())
	return reflect.NewAt(v.Type(), ptr).Elem()
}
