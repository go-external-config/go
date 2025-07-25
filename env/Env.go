package env

import (
	"fmt"
	"reflect"
	"strings"
)

var environment *Environment

// ${component.db.host}
// #{'${component.db.user}:${component.db.pass}'}
func Value(expression string) any {
	return EnvironmentInstance().ResolveRequiredPlaceholders(expression)
}

func ConfigurationProperties[T any](prefix string, target T) T {
	// var zero T
	t := reflect.TypeOf(target)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		value := Value(fmt.Sprintf("${%s.%s}", prefix, f.Name))
		fmt.Printf("%s[%s] = %s\n", f.Name, f.Type, value)
		v := reflect.ValueOf(&target).Elem()
		v.FieldByName(f.Name).Set(reflect.ValueOf(value))
	}
	return target
}

// last wins
func ActiveProfiles() []string {
	return EnvironmentInstance().activeProfiles
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

func EnvironmentInstance() *Environment {
	if environment == nil {
		environment = newEnvironment("")
	}
	return environment
}
