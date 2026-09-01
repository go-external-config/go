package env

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"unicode"

	"github.com/go-errr/go/err"
	"github.com/go-external-config/go/str"
	"github.com/go-jang/go/lang"
	refl "github.com/go-jang/go/lang/reflect"
	"github.com/go-jang/go/util/concurrent"
)

const ValueTag = "value"

func Property(key string) string {
	return Instance().property(key)
}

// Expression to evaluate against environment properties
//
//	require.Equal(t, "value", env.Value[string]("${key:default}"))
//	require.Equal(t, []string{"host1", "host2", "host3"}, env.Value[[]string]("#{split('${servers}', ',')}"))
func Value[T any](expression string) T {
	return convertAs[T](ResolveRequiredPlaceholders(expression))
}

func ResolvePlaceholders(expression string) any {
	return Instance().resolvePlaceholders(expression)
}

func ResolveRequiredPlaceholders(expression string) any {
	return Instance().resolveRequiredPlaceholders(expression)
}

// Binds properties with the given prefix to the target struct using field names
func ConfigurationProperties[T any](prefix string, target *T) *T {
	targetType := reflect.TypeOf(target).Elem()
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
		value := ResolveRequiredPlaceholders(rawValue.Value())
		targetFieldValue := targetValue.FieldByName(reflectField.Name)
		converted := convertAsType(value, targetFieldValue.Type())
		refl.Settable(targetFieldValue).Set(reflect.ValueOf(converted))
	}
	return target
}

// Binds properties to the target struct using field tags.
func BindProperties[T any](target *T) *T {
	BindPropertiesAny(target)
	return target
}

// Binds properties to the target struct using field tags.
func BindPropertiesAny(target any) any {
	refl.ForEachTaggedField(target, ValueTag, func(field refl.Field) {
		defer err.Catch(func(e any) {
			panic(err.NewRuntimeExceptionFrom(fmt.Sprintf("Cannot bind configuration value '%s' to field '%s'", field.TagValue, field.Field.Name), e))
		})
		value := ResolveRequiredPlaceholders(field.TagValue)
		converted := convertAsType(value, field.Type)
		field.Value.Set(reflect.ValueOf(converted))
	})

	return target
}

// last wins
func ActiveProfiles() []string {
	return Instance().activeProfiles()
}

// Determine whether one or more of the given profiles is active.
//
// If a profile begins with '!' the logic is inverted, meaning this method will return true if the given profile is not active.
// For example, env.MatchesProfiles("p1", "!p2") will return true if profile 'p1' is active or 'p2' is not active.
// A compound expression allows for more complicated profile logic to be expressed, for example "production & cloud".
func MatchesProfiles(profiles ...string) bool {
	return Instance().matchesProfiles(profiles...)
}

// Bootstrap new environment with profiles listed, last wins.
// Do nothing if the same profiles are already set in the specified order.
func SetActiveProfiles(profiles string) {
	concurrent.Synchronized(&environmentMu, func() {
		if environment != nil && slices.Equal(environment.activeProfiles(), splitProfiles(profiles)) {
			return
		}

		environment = newEnvironment(profiles)
	})
}

func splitProfiles(profiles string) []string {
	profiles = strings.TrimSpace(profiles)
	if profiles == "" {
		return nil
	}
	return profileSeparator.Split(profiles, -1)
}

func PropertySources() []PropertySource {
	return Instance().propertySources()
}

// Register custom property source to implement additional logic for properties processing, like property=base64:dGVzdAo=.
// See Base64PropertySource (available by default) and RsaPropertySource
func RegisterPropertySource(source PropertySource) {
	concurrent.Synchronized(&environmentMu, func() {
		registeredPropertySources = append(registeredPropertySources, source)
		if environment != nil {
			environment.addPropertySource(source)
		}
	})
}

// Add custom context variables to be evaluated.
// See env.ExprProcessor for expressions and variables available by default.
//
//	env.SetContextVariable("runtime", map[string]any{
//		"NumCPU": runtime.NumCPU(),
//	})
func SetContextVariable(key string, value any) {
	concurrent.Synchronized(&environmentMu, func() {
		contextVariables[key] = value
		if environment != nil {
			environment.setContextVariable(key, value)
		}
	})
}

func convertAs[T any](value any) T {
	return convertAsType(value, lang.TypeOf[T]()).(T)
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
			lang.Assert(val.Type().ConvertibleTo(t), "Cannot convert %s %v to %v", val.Type().Name(), value, t.Name())
			return val.Convert(t).Interface()
		}
	}
}
