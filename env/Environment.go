package env

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-external-config/go/lang"
	"github.com/go-external-config/go/util/collection"
	"github.com/go-external-config/go/util/files"
	"github.com/go-external-config/go/util/optional"
	"github.com/go-external-config/go/util/regex"
	"github.com/go-external-config/go/util/str"
	"github.com/go-external-config/go/util/text"
)

var locationPattern = regexp.MustCompile(`(?P<location>.+)\[(?P<fantomExt>\.[\w]+)\]`)
var envVarCanonicalFormTranslationRule = map[rune]rune{
	'.': '_',
	'[': '_',
	']': '_',
	'-': 0, // delete
}

type Environment struct {
	activeProfiles        []string
	paramsPropertySource  *MapPropertySource
	environPropertySource *MapPropertySource
	propertySources       []PropertySource
	exprProcessor         *ExprProcessor
}

func newEnvironment(activeProfiles string) *Environment {
	environment := Environment{
		activeProfiles:  []string{"default"},
		propertySources: make([]PropertySource, 0),
		exprProcessor:   ExprProcessorOf(true)}

	environment.loadEnvironmentVariables()
	environment.loadApplicationParameters()
	environment.loadApplicationConfiguration(activeProfiles)
	environment.AddPropertySource(NewRandomValuePropertySource())
	environment.AddPropertySource(NewBase64PropertySource(&environment))
	return &environment
}

func (e *Environment) Property(key string) string {
	return fmt.Sprint(e.ResolveRequiredPlaceholders(e.lookupRawProperty(key).
		OrElsePanic("No value present for %s", key)))
}

func (e *Environment) lookupRawProperty(key string) *optional.Optional[string] {
	if e.paramsPropertySource.HasProperty(key) {
		return optional.OfValue(e.paramsPropertySource.Property(key))
	} else if e.environPropertySource.HasProperty(key) {
		return optional.OfValue(e.environPropertySource.Property(key))
	} else if envCanonical := e.envVarCanonicalForm(key); e.environPropertySource.HasProperty(envCanonical) {
		return optional.OfValue(e.environPropertySource.Property(envCanonical))
	} else {
		for i := len(e.propertySources) - 1; i >= 0; i-- {
			if e.propertySources[i].HasProperty(key) {
				return optional.OfValue(e.propertySources[i].Property(key))
			}
		}
	}
	return optional.OfEmpty[string]()
}

func (e *Environment) ResolveRequiredPlaceholders(expression string) any {
	return e.exprProcessor.Process(expression)
}

// Determine whether one or more of the given profiles is active.
//
// If a profile begins with '!' the logic is inverted, meaning this method will return true if the given profile is not active.
// For example, env.MatchesProfiles("p1", "!p2") will return true if profile 'p1' is active or 'p2' is not active.
// A compound expression allows for more complicated profile logic to be expressed, for example "production & cloud".
func (e *Environment) MatchesProfiles(profiles ...string) bool {
	if len(profiles) == 0 {
		return true
	}
	activeProfiles := collection.SliceToSet(e.activeProfiles)
	processor := text.PatternProcessorOf("(?P<word>\\w+)|(?P<sign>\\W)")
	processor.OverrideResolve(func(match *regex.Match,
		super func(*regex.Match) any) any {
		word := match.NamedGroup("word")
		sign := match.NamedGroup("sign")

		if word.Present() {
			if _, found := activeProfiles[word.Value()]; found {
				return true
			} else {
				return false
			}
		}
		switch sign.Value() {
		case "&":
			return "&&"
		case "|":
			return "||"
		default:
			return match.Expr()
		}
	})
	for _, profile := range profiles {
		if Value[bool](fmt.Sprintf("#{%v}", processor.ProcessRecursive(profile, false))) {
			return true
		}
	}
	return false
}

// last wins
func (e *Environment) ActiveProfiles() []string {
	return e.activeProfiles
}

// first wins
func (e *Environment) PropertySources() []PropertySource {
	return collection.ReverseSlice(e.propertySources)
}

// ACTIVE_PROFILES=dev,hsqldb
func (e *Environment) loadEnvironmentVariables() {
	environ := MapPropertySourceOf("Environment variables")
	pattern := regexp.MustCompile(`(?P<key>[^=\s]+)=(?P<value>.*)`)
	for _, keyValue := range os.Environ() {
		for _, m := range pattern.FindAllStringSubmatchIndex(keyValue, -1) {
			match := regex.MatchOf(pattern, keyValue, m)
			environ.SetProperty(match.NamedGroup("key").Value(), match.NamedGroup("value").Value())
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
			match := regex.MatchOf(pattern, keyValue, m)
			params.SetProperty(match.NamedGroup("key").Value(), match.NamedGroup("value").Value())
		}
	}
	e.paramsPropertySource = params
}

// last wins
// application.yaml
// application-<profile>.yaml
func (e *Environment) loadApplicationConfiguration(bootstrapProfiles string) {
	activeProfiles := lang.FirstNonEmpty(e.paramsPropertySource.properties["profiles.active"], e.environPropertySource.properties["PROFILES_ACTIVE"], bootstrapProfiles)
	e.activeProfiles = lang.If(len(activeProfiles) == 0, e.activeProfiles, append(e.activeProfiles, strings.Split(activeProfiles, ",")...))
	configName := lang.FirstNonEmpty(e.paramsPropertySource.properties["config.name"], e.environPropertySource.properties["CONFIG_NAME"], "application")
	defaultLocation := "./,./config/"
	additionalLocation := lang.FirstNonEmpty(e.paramsPropertySource.properties["config.additional-location"], e.environPropertySource.properties["CONFIG_ADDITIONALLOCATION"])
	extendedDefaultLocation := lang.If(len(additionalLocation) == 0, defaultLocation, defaultLocation+","+additionalLocation)
	configLocation := lang.FirstNonEmpty(e.paramsPropertySource.properties["config.location"], e.environPropertySource.properties["CONFIG_LOCATION"])
	extendedConfigLocation := lang.If(len(additionalLocation) == 0, configLocation, additionalLocation+","+configLocation)
	resolvedConfigLocation := lang.If(len(configLocation) == 0, extendedDefaultLocation, extendedConfigLocation)

	for _, location := range strings.Split(resolvedConfigLocation, ",") {
		for i := 0; i < len(e.activeProfiles); i++ {
			for _, locationGroup := range strings.Split(location, ";") {
				e.loadConfiguration(locationGroup, configName, e.activeProfiles[i])
			}
		}
	}
}

func (e *Environment) loadConfiguration(location, name, profile string) {
	location = filepath.ToSlash(location)
	var fantomExt string
	for _, m := range locationPattern.FindAllStringSubmatchIndex(location, -1) {
		match := regex.MatchOf(locationPattern, location, m)
		location = match.NamedGroup("location").Value()
		fantomExt = match.NamedGroup("fantomExt").Value()
	}

	if strings.HasSuffix(location, "/") {
		e.loadFile(files.RelativePath(location, lang.If(profile == "default", name+".yml", name+"-"+profile+".yml")), fantomExt)
		e.loadFile(files.RelativePath(location, lang.If(profile == "default", name+".yaml", name+"-"+profile+".yaml")), fantomExt)
		e.loadFile(files.RelativePath(location, lang.If(profile == "default", name+".properties", name+"-"+profile+".properties")), fantomExt)
	} else if len(fantomExt) > 0 {
		e.loadFile(lang.If(profile == "default", location, location+"-"+profile), fantomExt)
	} else {
		ext := filepath.Ext(location)
		e.loadFile(lang.If(profile == "default", location, location[:len(location)-len(ext)]+"-"+profile+ext), fantomExt)
	}
}

func (e *Environment) loadFile(path, fantomExt string) {
	if !files.Exists(path) {
		return
	}
	var result PropertySource
	slog.Info(fmt.Sprintf("%T: loading properties from %s", *e, path))
	ext := lang.FirstNonEmpty(fantomExt, filepath.Ext(path))
	lang.AssertState(len(ext) != 0, "Cannot load from location %s. If location supposed to be a directory use '/' at the end. Otherwise provide extension hint in square brackets like [.properties] to derive property source type", path)
	file := optional.OfCommaErr(os.Open(path)).OrElsePanic("Cannot open file %s", path)
	defer file.Close()
	content := string(optional.OfCommaErr(io.ReadAll(file)).OrElsePanic("Cannot read from %s", path))
	switch ext {
	case ".properties":
		result = NewPropertiesPropertySource(path, content)
	case ".yaml", ".yml":
		result = NewYamlPropertySource(path, content)
	default:
		panic(fmt.Sprintf("Cannot load from %s as %s file type is not supported. Use extension hint in square brackets like .env[.properties] to derive property source type", path, ext))
	}
	e.propertySources = append(e.propertySources, result)
	if result.HasProperty("profiles.active") && len(e.activeProfiles) == 1 && e.activeProfiles[0] == "default" {
		e.activeProfiles = append(e.activeProfiles, strings.Split(result.Property("profiles.active"), ",")...)
	}
	if result.HasProperty("config.import") {
		for _, location := range strings.Split(result.Property("config.import"), ",") {
			e.loadImport(path, location)
		}
	}
}

func (e *Environment) loadImport(path, location string) {
	var fantomExt string
	for _, m := range locationPattern.FindAllStringSubmatchIndex(location, -1) {
		match := regex.MatchOf(locationPattern, location, m)
		location = match.NamedGroup("location").Value()
		fantomExt = match.NamedGroup("fantomExt").Value()
	}
	location = filepath.ToSlash(location)
	lang.AssertState(!strings.HasSuffix(location, "/"), "Cannot load from location %s defined in %s. Directory import is not supported", location, path)
	e.loadFile(files.RelativePath(path, location), fantomExt)
}

func (e *Environment) envVarCanonicalForm(key string) string {
	return strings.ToUpper(str.ReplaceChars(key, envVarCanonicalFormTranslationRule))
}

// Add custom property source to implement additional logic for properties processing, like property=base64:dGVzdAo=.
// See Base64PropertySource (available by default) and RsaPropertySource
//
//	var environment = env.Instance().AddPropertySource(env.NewRsaPropertySource())
func (e *Environment) AddPropertySource(source PropertySource) *Environment {
	e.propertySources = append(e.propertySources, source)
	return e
}
