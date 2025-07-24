package env

import "strings"

var environment *Environment

// ${component.db.host}
// #{'${component.db.user}:${component.db.pass}'}
func Value(expression string) string {
	return GetEnvironment().ResolveRequiredPlaceholders(expression)
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
