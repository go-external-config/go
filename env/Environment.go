package env

type Environment struct {
	activeProfiles  []string
	propertySources map[string]PropertySource
}

func EnvironmentOf(activeProfiles ...string) *Environment {
	return &Environment{
		activeProfiles:  activeProfiles,
		propertySources: make(map[string]PropertySource)}
}

// Determine whether one or more of the given profiles is active.
//
// If a profile begins with '!' the logic is inverted, meaning this method will return true if the given profile is not active.
// For example, env.MatchesProfiles("p1", "!p2") will return true if profile 'p1' is active or 'p2' is not active.
func (e *Environment) MatchesProfiles(profiles ...string) {}

func (e *Environment) GetActiveProfiles() []string {
	return e.activeProfiles
}
