package env

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
	"unsafe"

	"github.com/go-external-config/go/lang"
	"github.com/go-external-config/go/util/str"
)

var environment *Environment

// Expression to evaluate against environment properties
//
//	require.Equal(t, "value", env.Value[string]("${key:default}"))
//	require.Equal(t, []string{"host1", "host2", "host3"}, env.Value[[]string]("#{split('${servers}', ',')}"))
func Value[T any](expression string) T {
	return convertAs[T](Instance().ResolveRequiredPlaceholders(expression))
}

func ConfigurationProperties[T any](prefix string, target *T) *T {
	targetType := reflect.TypeOf(*target)
	targetValue := reflect.ValueOf(target).Elem()
	for i := 0; i < targetType.NumField(); i++ {
		reflectField := targetType.Field(i)
		rawValue := Instance().lookupRawProperty(fmt.Sprintf("%s.%s", prefix, reflectField.Name))
		if !rawValue.Present() && unicode.IsUpper(rune(reflectField.Name[0])) {
			decapitalizedName := strings.ToLower(reflectField.Name[:1]) + reflectField.Name[1:]
			rawValue = Instance().lookupRawProperty(fmt.Sprintf("%s.%s", prefix, decapitalizedName))
		}
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

// Determine whether one or more of the given profiles is active.
//
// If a profile begins with '!' the logic is inverted, meaning this method will return true if the given profile is not active.
// For example, env.MatchesProfiles("p1", "!p2") will return true if profile 'p1' is active or 'p2' is not active.
// A compound expression allows for more complicated profile logic to be expressed, for example "production & cloud".
func MatchesProfiles(profiles ...string) bool {
	return Instance().MatchesProfiles(profiles...)
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
		switch v := value.(type) {
		case string:
			return v
		default:
			return fmt.Sprint(v)
		}
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
