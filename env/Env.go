package env

var environment *Environment = newEnvironment()

func Value(expression string) string {
	return environment.ResolveRequiredPlaceholders(expression)
}

func ActiveProfiles() []string {
	return environment.activeProfiles
}

func ActiveteProfiles(profiles string) {
	environment.loadApplicationConfiguration(profiles)
}

func GetEnvironment() *Environment {
	return environment
}
