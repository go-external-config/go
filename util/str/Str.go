package str

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/go-external-config/go/lang"
)

func Parse[T any](value string) T {
	var zero T
	t := reflect.TypeOf(zero)
	return ParseOfType(value, t).(T)
}

func ParseOfType(value string, t reflect.Type) any {
	errMsg := "Failed to parse '%s' as %s\nCaused by: %s"
	switch t.Kind() {
	case reflect.Int:
		v, err := strconv.Atoi(value)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return reflect.ValueOf(v).Convert(t).Interface()
	case reflect.Int8:
		v, err := strconv.ParseInt(value, 10, 8)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return reflect.ValueOf(int8(v)).Convert(t).Interface()
	case reflect.Int16:
		v, err := strconv.ParseInt(value, 10, 16)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return reflect.ValueOf(int16(v)).Convert(t).Interface()
	case reflect.Int32:
		v, err := strconv.ParseInt(value, 10, 32)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return reflect.ValueOf(int32(v)).Convert(t).Interface()
	case reflect.Int64:
		v, err := strconv.ParseInt(value, 10, 64)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return reflect.ValueOf(v).Convert(t).Interface()
	case reflect.Uint:
		v, err := strconv.ParseUint(value, 10, 0)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return reflect.ValueOf(uint(v)).Convert(t).Interface()
	case reflect.Uint8:
		v, err := strconv.ParseUint(value, 10, 8)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return reflect.ValueOf(uint8(v)).Convert(t).Interface()
	case reflect.Uint16:
		v, err := strconv.ParseUint(value, 10, 16)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return reflect.ValueOf(uint16(v)).Convert(t).Interface()
	case reflect.Uint32:
		v, err := strconv.ParseUint(value, 10, 32)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return reflect.ValueOf(uint32(v)).Convert(t).Interface()
	case reflect.Uint64:
		v, err := strconv.ParseUint(value, 10, 64)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return reflect.ValueOf(v).Convert(t).Interface()
	case reflect.Float32:
		v, err := strconv.ParseFloat(value, 32)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return reflect.ValueOf(float32(v)).Convert(t).Interface()
	case reflect.Float64:
		v, err := strconv.ParseFloat(value, 64)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return reflect.ValueOf(v).Convert(t).Interface()
	case reflect.Bool:
		v, err := strconv.ParseBool(value)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return reflect.ValueOf(v).Convert(t).Interface()
	case reflect.String:
		return reflect.ValueOf(value).Convert(t).Interface()
	default:
		panic(fmt.Sprintf("Unsupported type: %s", t))
	}
}
