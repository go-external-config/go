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
		return v
	case reflect.Int8:
		v, err := strconv.ParseInt(value, 10, 8)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return int8(v)
	case reflect.Int16:
		v, err := strconv.ParseInt(value, 10, 16)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return int16(v)
	case reflect.Int32:
		v, err := strconv.ParseInt(value, 10, 32)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return int32(v)
	case reflect.Int64:
		v, err := strconv.ParseInt(value, 10, 64)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return v
	case reflect.Uint:
		v, err := strconv.ParseUint(value, 10, 0)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return uint(v)
	case reflect.Uint8:
		v, err := strconv.ParseUint(value, 10, 8)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return uint8(v)
	case reflect.Uint16:
		v, err := strconv.ParseUint(value, 10, 16)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return uint16(v)
	case reflect.Uint32:
		v, err := strconv.ParseUint(value, 10, 32)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return uint32(v)
	case reflect.Uint64:
		v, err := strconv.ParseUint(value, 10, 64)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return v
	case reflect.Float32:
		v, err := strconv.ParseFloat(value, 32)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return float32(v)
	case reflect.Float64:
		v, err := strconv.ParseFloat(value, 64)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return v
	case reflect.Bool:
		v, err := strconv.ParseBool(value)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return v
	case reflect.String:
		return value
	default:
		panic(fmt.Sprintf("Unsupported type: %s", t))
	}
}
