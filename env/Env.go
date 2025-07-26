package env

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"github.com/go-external-config/go/lang"
	"github.com/go-external-config/go/util/str"
)

var environment *Environment

// Expression to evaluate against environment properties
//
//	require.Equal(t, "value", env.Value("${key}"))
//	require.Equal(t, []string{"host1", "host2", "host3"}, env.Value("#{split('${servers}', ',')}"))
func Value(expression string) any {
	return Instance().ResolveRequiredPlaceholders(expression)
}

func ValueAs[T any](expression string) T {
	return convertAs[T](Value(expression))
}

// Lookup property value, evaluate any expressions if any.
// Properties are strings unless value is an expression which evaluates to any type
//
//	require.Equal(t, "value", env.Property("key"))
//	require.Equal(t, "value", env.Value("${key}"))
//	require.Equal(t, []string{"host1", "host2", "host3"}, env.Value("#{split('${servers}', ',')}"))
func Property(prop string) any {
	return Instance().Property(prop)
}

func PropertyAs[T any](expression string) T {
	return convertAs[T](Property(expression))
}

func ConfigurationProperties[T any](prefix string, target *T) *T {
	targetType := reflect.TypeOf(*target)
	targetValue := reflect.ValueOf(target).Elem()
	for i := 0; i < targetType.NumField(); i++ {
		reflectField := targetType.Field(i)
		rawValue := Instance().lookupRawProperty(fmt.Sprintf("%s.%s", prefix, reflectField.Name))
		if !rawValue.Present() {
			continue
		}
		value := Instance().ResolveRequiredPlaceholders(rawValue.Value())
		targetFieldValue := targetValue.FieldByName(reflectField.Name)
		converted := convertAsType(value, targetFieldValue.Type())
		if targetFieldValue.CanSet() {
			targetFieldValue.Set(reflect.ValueOf(converted))
		} else {
			ptr := unsafe.Pointer(targetFieldValue.UnsafeAddr())
			settableField := reflect.NewAt(targetFieldValue.Type(), ptr).Elem()
			settableField.Set(reflect.ValueOf(converted))
		}
	}
	return target
}

// last wins
func ActiveProfiles() []string {
	return Instance().activeProfiles
}

// Bootstrap new environment with profiles listed, last wins.
// Do nothing if very profiles are already set in the specified order.
// Reload environment if empty value provided
func SetActiveProfiles(profiles string) {
	if environment != nil && "default,"+profiles == strings.Join(environment.ActiveProfiles(), ",") {
		return
	}

	previous := environment
	environment = newEnvironment(profiles)

	// keep custom property preprocessors
	if previous != nil {
		for _, source := range previous.propertySources {
			if source.Properties() == nil {
				environment.AddPropertySource(source)
			}
		}
	}
}

func Instance() *Environment {
	if environment == nil {
		environment = newEnvironment("")
	}
	return environment
}

func convertAs[T any](value any) T {
	var zero T
	return convertAsType(value, reflect.TypeOf(zero)).(T)
}

func convertAsType(value any, t reflect.Type) any {
	switch t.Kind() {
	case reflect.String:
		return fmt.Sprintf("%v", value)
	default:
		switch v := value.(type) {
		case string:
			return str.ParseOfType(v, t)
		default:
			val := reflect.ValueOf(value)
			lang.AssertState(val.Type().ConvertibleTo(t), "Cannot convert %s %v to %v", val.Type().Name(), value, t.Name())
			return val.Convert(t).Interface()
		}
	}
}
