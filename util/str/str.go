package str

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/go-errr/go/err"
	"github.com/go-external-config/go/lang"
	"github.com/go-external-config/go/util/optional"
)

func Parse[T any](value string) T {
	t := lang.TypeOf[T]()
	return ParseOfType(value, t).(T)
}

func ParseOfType(value string, t reflect.Type) any {
	errMsg := "Failed to parse '%s' as %s\nCaused by: %s"
	switch t.Kind() {
	case reflect.Int:
		v, e := strconv.Atoi(value)
		if e != nil {
			panic(err.NewNumberFormatException(fmt.Sprintf(errMsg, value, t, e)))
		}
		return reflect.ValueOf(v).Convert(t).Interface()
	case reflect.Int8:
		v, e := strconv.ParseInt(value, 10, 8)
		if e != nil {
			panic(err.NewNumberFormatException(fmt.Sprintf(errMsg, value, t, e)))
		}
		return reflect.ValueOf(int8(v)).Convert(t).Interface()
	case reflect.Int16:
		v, e := strconv.ParseInt(value, 10, 16)
		if e != nil {
			panic(err.NewNumberFormatException(fmt.Sprintf(errMsg, value, t, e)))
		}
		return reflect.ValueOf(int16(v)).Convert(t).Interface()
	case reflect.Int32:
		v, e := strconv.ParseInt(value, 10, 32)
		if e != nil {
			panic(err.NewNumberFormatException(fmt.Sprintf(errMsg, value, t, e)))
		}
		return reflect.ValueOf(int32(v)).Convert(t).Interface()
	case reflect.Int64:
		v, e := strconv.ParseInt(value, 10, 64)
		if e != nil {
			panic(err.NewNumberFormatException(fmt.Sprintf(errMsg, value, t, e)))
		}
		return reflect.ValueOf(v).Convert(t).Interface()
	case reflect.Uint:
		v, e := strconv.ParseUint(value, 10, 0)
		if e != nil {
			panic(err.NewNumberFormatException(fmt.Sprintf(errMsg, value, t, e)))
		}
		return reflect.ValueOf(uint(v)).Convert(t).Interface()
	case reflect.Uint8:
		v, e := strconv.ParseUint(value, 10, 8)
		if e != nil {
			panic(err.NewNumberFormatException(fmt.Sprintf(errMsg, value, t, e)))
		}
		return reflect.ValueOf(uint8(v)).Convert(t).Interface()
	case reflect.Uint16:
		v, e := strconv.ParseUint(value, 10, 16)
		if e != nil {
			panic(err.NewNumberFormatException(fmt.Sprintf(errMsg, value, t, e)))
		}
		return reflect.ValueOf(uint16(v)).Convert(t).Interface()
	case reflect.Uint32:
		v, e := strconv.ParseUint(value, 10, 32)
		if e != nil {
			panic(err.NewNumberFormatException(fmt.Sprintf(errMsg, value, t, e)))
		}
		return reflect.ValueOf(uint32(v)).Convert(t).Interface()
	case reflect.Uint64:
		v, e := strconv.ParseUint(value, 10, 64)
		if e != nil {
			panic(err.NewNumberFormatException(fmt.Sprintf(errMsg, value, t, e)))
		}
		return reflect.ValueOf(v).Convert(t).Interface()
	case reflect.Float32:
		v, e := strconv.ParseFloat(value, 32)
		if e != nil {
			panic(err.NewNumberFormatException(fmt.Sprintf(errMsg, value, t, e)))
		}
		return reflect.ValueOf(float32(v)).Convert(t).Interface()
	case reflect.Float64:
		v, e := strconv.ParseFloat(value, 64)
		if e != nil {
			panic(err.NewNumberFormatException(fmt.Sprintf(errMsg, value, t, e)))
		}
		return reflect.ValueOf(v).Convert(t).Interface()
	case reflect.Bool:
		v, e := strconv.ParseBool(value)
		if e != nil {
			panic(err.NewNumberFormatException(fmt.Sprintf(errMsg, value, t, e)))
		}
		return reflect.ValueOf(v).Convert(t).Interface()
	case reflect.String:
		return reflect.ValueOf(value).Convert(t).Interface()
	default:
		panic(err.NewIllegalArgumentException(fmt.Sprintf("Unsupported type: %s", t)))
	}
}

func ReplaceChars(str string, rules map[rune]rune) string {
	var builder strings.Builder
	builder.Grow(len(str))
	for _, rule := range str {
		if replacement, ok := rules[rule]; ok {
			if replacement != 0 { // 0 means 'delete'
				builder.WriteRune(replacement)
			}
		} else {
			builder.WriteRune(unicode.ToUpper(rule))
		}
	}
	return builder.String()
}

func Join(delim string, values ...any) string {
	if len(values) == 0 {
		return ""
	}
	var delimRequired = false
	var b strings.Builder
	for _, value := range values {
		if delimRequired {
			b.WriteString(delim)
		}
		delimRequired = optional.OfCommaErr(b.WriteString(fmt.Sprint(value))).OrElsePanic("Cannot write %v", value) > 0
	}
	return b.String()
}
