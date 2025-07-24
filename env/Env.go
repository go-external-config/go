package env

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/madamovych/go/lang"
)

var environment *Environment

// ${component.db.host}
// #{'${component.db.user}:${component.db.pass}'}
func Value[T any](expression string) T {

	value := GetEnvironment().ResolveRequiredPlaceholders(expression)

	var zero T
	t := reflect.TypeOf(zero).String()
	errMsg := "Failed to parse '%s' as %s\nCaused by: %s"
	switch any(zero).(type) {
	case int:
		v, err := strconv.Atoi(value)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return any(v).(T)
	case int8:
		v, err := strconv.ParseInt(value, 10, 8)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return any(int8(v)).(T)
	case int16:
		v, err := strconv.ParseInt(value, 10, 16)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return any(int16(v)).(T)
	case int32:
		v, err := strconv.ParseInt(value, 10, 32)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return any(int32(v)).(T)
	case int64:
		v, err := strconv.ParseInt(value, 10, 64)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return any(v).(T)
	case uint:
		v, err := strconv.ParseUint(value, 10, 0)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return any(uint(v)).(T)
	case uint8:
		v, err := strconv.ParseUint(value, 10, 8)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return any(uint8(v)).(T)
	case uint16:
		v, err := strconv.ParseUint(value, 10, 16)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return any(uint16(v)).(T)
	case uint32:
		v, err := strconv.ParseUint(value, 10, 32)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return any(uint32(v)).(T)
	case uint64:
		v, err := strconv.ParseUint(value, 10, 64)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return any(v).(T)
	case float32:
		v, err := strconv.ParseFloat(value, 32)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return any(float32(v)).(T)
	case float64:
		v, err := strconv.ParseFloat(value, 64)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return any(v).(T)
	case bool:
		v, err := strconv.ParseBool(value)
		lang.AssertState(err == nil, errMsg, value, t, err)
		return any(v).(T)
	case string:
		return any(value).(T)
	default:
		panic(fmt.Sprintf("Unsupported type: %s", t))
	}
}

// last wins
func ActiveProfiles() []string {
	return GetEnvironment().activeProfiles
}

// bootstrap new environment with profiles listed, last wins. Do nothing if very profiles are already set
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

func GetEnvironment() *Environment {
	if environment == nil {
		environment = newEnvironment("")
	}
	return environment
}
