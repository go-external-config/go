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
	"github.com/madamovych/go/util"
	"github.com/madamovych/go/util/collection"
	"github.com/madamovych/go/util/regex"
	"github.com/madamovych/go/util/text"
)

type Environment struct {
	activeProfiles        []string
	paramsPropertySource  *MapPropertySource
	environPropertySource *MapPropertySource
	propertySources       []PropertySource
	resourceLoader        *io.ResourceLoader
	exprProcessor         *text.ExprProcessor
}

func newEnvironment() *Environment {
	environment := Environment{
		activeProfiles:  []string{"default"},
		propertySources: make([]PropertySource, 0),
		resourceLoader:  io.NewResourceLoader(),
		exprProcessor:   text.ExprProcessorOf(true)}

	environment.loadEnvironmentVariables()
	environment.loadApplicationParameters()
	environment.loadApplicationConfiguration("")
	return &environment
}

func (e *Environment) Property(key string) string {
	if e.paramsPropertySource.HasProperty(key) {
		return e.ResolveRequiredPlaceholders(e.paramsPropertySource.Property(key))
	} else if e.environPropertySource.HasProperty(key) {
		return e.ResolveRequiredPlaceholders(e.environPropertySource.Property(key))
	} else {
		for i := len(e.propertySources) - 1; i >= 0; i-- {
			if e.propertySources[i].HasProperty(key) {
				return e.ResolveRequiredPlaceholders(e.propertySources[i].Property(key))
			}
		}
		panic("No value present for " + key)
	}
}

func (e *Environment) ResolveRequiredPlaceholders(expression string) string {
	return e.exprProcessor.Process(expression)
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
func (e *Environment) loadEnvironmentVariables() {
	environ := MapPropertySourceOf("Environment variables")
	pattern := regexp.MustCompile(`(?P<key>[^=\s]+)=(?P<value>.*)`)
	for _, keyValue := range os.Environ() {
		for _, m := range pattern.FindAllStringSubmatchIndex(keyValue, -1) {
			match := regex.RegexpMatchOf(pattern, keyValue, m)
			key := match.NamedGroup("key").Value()
			value := match.NamedGroup("value").Value()
			environ.SetProperty(key, value)
			e.exprProcessor.Define(key, value)
		}
	}
	e.environPropertySource = environ
}

// --active.profiles=dev,hsqldb
func (e *Environment) loadApplicationParameters() {
	params := MapPropertySourceOf("Application parameters")
	pattern := regexp.MustCompile(`--?(?P<key>[^=\s]+)\s*=?(?P<value>.*)`)
	for _, keyValue := range os.Args[1:] {
		for _, m := range pattern.FindAllStringSubmatchIndex(keyValue, -1) {
			match := regex.RegexpMatchOf(pattern, keyValue, m)
			key := match.NamedGroup("key").Value()
			value := match.NamedGroup("value").Value()
			params.SetProperty(key, value)
			e.exprProcessor.Define(key, value)
		}
	}
	e.paramsPropertySource = params
}

// application.yaml
// application-<profile>.yaml
func (e *Environment) loadApplicationConfiguration(additionalProfiles string) {
	bootstrap := len(additionalProfiles) == 0
	missingProfiles := collection.SubtractSlice(strings.Split(additionalProfiles, ","), e.activeProfiles)
	if bootstrap {
		var activeProfiles string = lang.FirstNonEmpty(e.paramsPropertySource.properties["profiles.active"], e.environPropertySource.properties["PROFILES_ACTIVE"])
		e.activeProfiles = lang.If(len(activeProfiles) == 0, e.activeProfiles, append(e.activeProfiles, strings.Split(activeProfiles, ",")...))
	} else {
		e.activeProfiles = lang.If(len(missingProfiles) == 0, e.activeProfiles, append(e.activeProfiles, missingProfiles...))
	}
	configName := lang.FirstNonEmpty(e.paramsPropertySource.properties["config.name"], e.environPropertySource.properties["CONFIG_NAME"], "application")
	configLocation := lang.FirstNonEmpty(e.paramsPropertySource.properties["config.location"], e.environPropertySource.properties["CONFIG_LOCATION"], "./")
	configAdditionalLocation := lang.FirstNonEmpty(e.paramsPropertySource.properties["config.additional-location"], e.environPropertySource.properties["CONFIG_ADDITIONAL_LOCATION"])
	configLocation = lang.If(len(configAdditionalLocation) == 0, configLocation, configLocation+";"+configAdditionalLocation)

	for _, location := range strings.Split(configLocation, ",") {
		for i := 0; i < len(lang.If(bootstrap, e.activeProfiles, missingProfiles)); i++ {
			for _, locationGroup := range strings.Split(location, ";") {
				e.loadConfiguration(locationGroup, configName, lang.If(bootstrap, e.activeProfiles, missingProfiles)[i])
			}
		}
	}
	e.loadEnvironmentVariables()
	e.loadApplicationParameters()
}

func (e *Environment) loadConfiguration(location, name, profile string) {
	locationPattern := regexp.MustCompile(`(?P<location>.+)\[(?P<fantomExt>\.[\w]+)\]`)
	var fantomExt string
	for _, m := range locationPattern.FindAllStringSubmatchIndex(location, -1) {
		match := regex.RegexpMatchOf(locationPattern, location, m)
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
	content := string(util.OptionalOfCommaErr(goio.ReadAll(resource.Reader())).OrElsePanic("Cannot read from %s", resource.URL().Path))
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
				e.activeProfiles = append(e.activeProfiles, strings.Split(value, ",")...)
			} else {
				continue
			}
		}
		e.exprProcessor.Define(key, value)
	}
}
