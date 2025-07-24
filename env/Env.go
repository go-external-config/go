package env

var environment *Environment = newEnvironment()

// ${component.db.host}
// #{'${component.db.user}:${component.db.pass}'}
func Value(expression string) string {
	return environment.ResolveRequiredPlaceholders(expression)
}

// last wins
func ActiveProfiles() []string {
	return environment.activeProfiles
}

// Activate additional profiles after bootstrap, last wins. Nothing happens if profile is already actived
func ActivateProfiles(profiles string) {
	environment.loadApplicationConfiguration(profiles)
}

func GetEnvironment() *Environment {
	return environment
}

func Refresh() {
	environment = newEnvironment()
}
