package env

import (
	"fmt"
	goio "io"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/madamovych/go/io"
	"github.com/madamovych/go/lang"
	"github.com/madamovych/go/text"
)

type Environment struct {
	activeProfiles        []string
	propertySources       []PropertySource
	applicationParameters map[string]string
	environmentVariables  map[string]string
	resourceLoader        *io.ResourceLoader
	exprProcessor         *text.ExprProcessor
}

func NewEnvironment() *Environment {
	environment := Environment{
		activeProfiles:  []string{"default"},
		propertySources: make([]PropertySource, 0),
		resourceLoader:  io.NewResourceLoader(),
		exprProcessor:   text.ExprProcessorOf(false)}

	environment.SystemEnvironmment()
	environment.ApplicationParameters()
	environment.ApplicationConfiguration()
	return &environment
}

// Determine whether one or more of the given profiles is active.
//
// If a profile begins with '!' the logic is inverted, meaning this method will return true if the given profile is not active.
// For example, env.MatchesProfiles("p1", "!p2") will return true if profile 'p1' is active or 'p2' is not active.
func (e *Environment) MatchesProfiles(profiles ...string) {}

func (e *Environment) ActiveProfiles() []string {
	return e.activeProfiles
}

func (e *Environment) PropertySources() []PropertySource {
	return e.propertySources
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
// application-<profile>.yaml
func (e *Environment) ApplicationConfiguration() {
	activeProfiles := lang.FirstNonEmpty(e.applicationParameters["profiles.active"], e.environmentVariables["PROFILES_ACTIVE"])
	e.activeProfiles = lang.If(len(activeProfiles) == 0, e.activeProfiles, append(e.activeProfiles, strings.Split(activeProfiles, ",")...))
	configName := lang.FirstNonEmpty(e.applicationParameters["config.name"], e.environmentVariables["CONFIG_NAME"], "application")
	configLocation := lang.FirstNonEmpty(e.applicationParameters["config.location"], e.environmentVariables["CONFIG_LOCATION"], "./")
	configAdditionalLocation := lang.FirstNonEmpty(e.applicationParameters["config.additional-location"], e.environmentVariables["CONFIG_ADDITIONAL_LOCATION"])
	configLocation = lang.If(len(configAdditionalLocation) == 0, configLocation, configLocation+";"+configAdditionalLocation)

	for _, location := range strings.Split(configLocation, ",") {
		for i := 0; i < len(e.activeProfiles); i++ {
			for _, locationGroup := range strings.Split(location, ";") {
				e.loadConfiguration(locationGroup, configName, e.activeProfiles[i])
			}
		}
	}

	environ := MapPropertySourceOf("Environment variables")
	for key, value := range e.environmentVariables {
		environ.SetProperty(key, string(value))
	}
	e.propertySources = append(e.propertySources, environ)

	params := MapPropertySourceOf("Application parameters")
	for key, value := range e.applicationParameters {
		params.SetProperty(key, string(value))
	}
	e.propertySources = append(e.propertySources, params)
}

func (e *Environment) loadConfiguration(location, name, profile string) {
	locationPattern := regexp.MustCompile(`(?P<location>.+)\[(?P<fantomExt>\.[\w]+)\]`)
	var fantomExt string
	for _, m := range locationPattern.FindAllStringSubmatchIndex(location, -1) {
		match := lang.RegexpMatchOf(locationPattern, location, m)
		location = match.NamedGroup("location").Value()
		fantomExt = match.NamedGroup("fantomExt").Value()
	}

	resource := e.resourceLoader.Resolve(location)
	if strings.HasSuffix(location, "/") {
		e.tryLoad(resource.CreateRelative(lang.If(profile == "default", name+".yml", name+"-"+profile+".yml")), fantomExt)
		e.tryLoad(resource.CreateRelative(lang.If(profile == "default", name+".yaml", name+"-"+profile+".yaml")), fantomExt)
		e.tryLoad(resource.CreateRelative(lang.If(profile == "default", name+".properties", name+"-"+profile+".properties")), fantomExt)
		name = "config/" + name
		e.tryLoad(resource.CreateRelative(lang.If(profile == "default", name+".yml", name+"-"+profile+".yml")), fantomExt)
		e.tryLoad(resource.CreateRelative(lang.If(profile == "default", name+".yaml", name+"-"+profile+".yaml")), fantomExt)
		e.tryLoad(resource.CreateRelative(lang.If(profile == "default", name+".properties", name+"-"+profile+".properties")), fantomExt)
	} else {
		e.tryLoad(lang.If(profile == "default", resource, resource.CreateRelative("-"+profile)), fantomExt)
	}
}

func (e *Environment) tryLoad(resource io.Resource, fantomExt string) {
	var result PropertySource
	if !resource.Exists() {
		return
	}
	slog.Info(fmt.Sprintf("Loading properties from %s", resource.URL().String()))
	ext := lang.FirstNonEmpty(filepath.Ext(resource.URL().Path), fantomExt)
	lang.AssertState(len(ext) != 0, "Cannot load from location %s. Either use '/' at the end if location supposed to be a directory or provide fantom extension like [.yaml] to derive property source type", resource.URL())
	content := string(lang.OptionalOfCommaErr(goio.ReadAll(resource.Reader())).OrElsePanic("Cannot read from %s", resource.URL().Path))
	switch ext {
	case ".properties":
		result = NewPropertiesPropertySource(resource.URL().Path, content)
	case ".yaml", ".yml":
		result = NewYamlPropertySource(resource.URL().Path, content)
	default:
		panic(fmt.Sprintf("Cannot load from %s as %s file types is not supported", resource.URL(), ext))
	}
	e.propertySources = append(e.propertySources, result)
	for key, value := range result.Properties() {
		if key == "active.profiles" {
			if len(e.activeProfiles) == 1 && e.activeProfiles[0] == "default" {
				e.activeProfiles = append(e.activeProfiles, strings.Split(value.(string), ",")...)
			} else {
				continue
			}
		}
		e.exprProcessor.Define(key, value)
	}
}
