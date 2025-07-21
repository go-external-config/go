package env

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/madamovych/go/lang"
)

type Environment struct {
	activeProfiles        []string
	propertySources       map[string]PropertySource
	applicationParameters map[string]string
	environmentVariables  map[string]string
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

func (e *Environment) ActiveProfiles() []string {
	return e.activeProfiles
}

// ACTIVE_PROFILES=dev,hsqldb
func (e *Environment) SystemEnvironmment() map[string]string {
	if e.environmentVariables != nil {
		return e.environmentVariables
	}
	properties := make(map[string]string)
	pattern := regexp.MustCompile(`(?P<key>[^=\s]+)=(?P<value>.*)`)
	for _, keyValue := range os.Environ() {
		for _, m := range pattern.FindAllStringSubmatchIndex(keyValue, -1) {
			match := lang.RegexpMatchOf(pattern, keyValue, m)
			properties[match.NamedGroup("key").Value()] = match.NamedGroup("value").Value()
		}
	}
	e.environmentVariables = properties
	return properties
}

// -flag -account=Stackoverflow
func (e *Environment) ApplicationParameters() map[string]string {
	if e.applicationParameters != nil {
		return e.applicationParameters
	}
	properties := make(map[string]string)
	pattern := regexp.MustCompile(`--?(?P<key>[^=\s]+)\s*=?(?P<value>.*)`)
	for _, keyValue := range os.Args[1:] {
		for _, m := range pattern.FindAllStringSubmatchIndex(keyValue, -1) {
			match := lang.RegexpMatchOf(pattern, keyValue, m)
			properties[match.NamedGroup("key").Value()] = match.NamedGroup("value").Value()
		}
	}
	e.applicationParameters = properties
	return properties
}

// application.yaml
func (e *Environment) ApplicationConfiguration() map[string]string {
	activeProfiles := lang.FirstNonEmpty(e.applicationParameters["profiles.active"], e.environmentVariables["PROFILES_ACTIVE"], "default")
	configName := lang.FirstNonEmpty(e.applicationParameters["config.name"], e.environmentVariables["CONFIG_NAME"], "application")
	configLocation := lang.FirstNonEmpty(e.applicationParameters["config.location"], e.environmentVariables["CONFIG_LOCATION"], "./")
	configAdditionalLocation := lang.FirstNonEmpty(e.applicationParameters["config.additional-location"], e.environmentVariables["CONFIG_ADDITIONAL_LOCATION"])
	configLocation = lang.If(len(configAdditionalLocation) == 0, configLocation, configLocation+";"+configAdditionalLocation)

	var properties map[string]any
	for _, locationGroup := range strings.Split(configLocation, ";") {
		for _, profile := range strings.Split(activeProfiles, ",") {
			for _, location := range strings.Split(locationGroup, ",") {
				e.loadConfiguration(location, configName, profile, properties)
			}
		}
	}

	return nil
}

func (e *Environment) loadConfiguration(location, name, profile string, properties map[string]any) {
	url := lang.OptionalOfCommaErr(url.Parse(location)).OrElsePanic("Cannot parse location %s", location)
	schema := lang.FirstNonEmpty(url.Scheme, "file")

	fmt.Printf("schema[%s], path[%s]\n", url.Scheme, url.Path)
	fmt.Printf("loading location[%s], name[%s], profile[%s]", location, name, profile)
}
