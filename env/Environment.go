package env

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"

	"github.com/go-errr/go/err"
	"github.com/go-external-config/go/files"
	"github.com/go-external-config/go/str"
	"github.com/go-jang/go/lang"
	"github.com/go-jang/go/util/collections"
	"github.com/go-jang/go/util/concurrent"
	"github.com/go-jang/go/util/objects"
	"github.com/go-jang/go/util/optional"
	"github.com/go-jang/go/util/regex"
	"github.com/go-jang/go/util/stream"
)

var defaultProfile = "default"
var profilesActiveEnv = "PROFILES_ACTIVE"
var profilesIncludeEnv = "PROFILES_INCLUDE"
var profilesActiveProp = "profiles.active"
var profilesIncludeProp = "profiles.include"
var profileSeparator = regexp.MustCompile(`\s*,\s*`)
var locationPattern = regexp.MustCompile(regex.NewPatternBuilder().Next(`{location:.+}\[{fantomExt:\.[\w]+}\]`).Build())
var envVarCanonicalFormTranslationRule = map[rune]rune{
	'.': '_',
	'[': '_',
	']': '_',
	'-': 0, // delete
}

var environment *Environment
var environmentMu sync.Mutex
var registeredPropertySources []PropertySource
var contextVariables = make(map[string]any)

type Environment struct {
	profiles              []string
	includedProfiles      []string
	paramsPropertySource  *MapPropertySource
	environPropertySource *MapPropertySource
	sources               []PropertySource
	exprProcessor         *ExprProcessor
	strictExprProcessor   *ExprProcessor
}

func Instance() *Environment {
	if environment == nil {
		concurrent.Synchronized(&environmentMu, func() {
			if environment == nil {
				environment = newEnvironment("")
			}
		})
	}
	return environment
}

func newEnvironment(activeProfiles string) *Environment {
	environment := Environment{
		profiles:            make([]string, 0),
		includedProfiles:    make([]string, 0),
		sources:             make([]PropertySource, 0),
		exprProcessor:       NewExprProcessor(false),
		strictExprProcessor: NewExprProcessor(true)}

	for key, value := range contextVariables {
		environment.setContextVariable(key, value)
	}

	environment.loadEnvironmentVariables()
	environment.loadApplicationParameters()
	environment.loadApplicationConfiguration(activeProfiles)
	environment.addPropertySource(NewRandomValuePropertySource())
	environment.addPropertySource(NewBase64PropertySource())
	environment.addPropertySource(NewCachedPropertySource())

	for _, source := range registeredPropertySources {
		environment.addPropertySource(source)
	}

	return &environment
}

func (this *Environment) property(key string) string {
	return fmt.Sprint(this.resolveRequiredPlaceholders(this.lookupRawProperty(key).
		OrElsePanic("No value present for %s", key)))
}

func (this *Environment) lookupRawProperty(key string) *optional.Optional[string] {
	if this.paramsPropertySource.HasProperty(key) {
		return optional.OfValue(this.paramsPropertySource.Property(key))
	} else if this.environPropertySource.HasProperty(key) {
		return optional.OfValue(this.environPropertySource.Property(key))
	} else if envCanonical := this.envVarCanonicalForm(key); this.environPropertySource.HasProperty(envCanonical) {
		return optional.OfValue(this.environPropertySource.Property(envCanonical))
	} else {
		for i := len(this.sources) - 1; i >= 0; i-- {
			if this.sources[i].HasProperty(key) {
				return optional.OfValue(this.sources[i].Property(key))
			}
		}
	}
	return optional.OfEmpty[string]()
}

func (this *Environment) resolvePlaceholders(expression string) any {
	return this.exprProcessor.Process(expression)
}

func (this *Environment) resolveRequiredPlaceholders(expression string) any {
	return this.strictExprProcessor.Process(expression)
}

func (this *Environment) matchesProfiles(profiles ...string) bool {
	if len(profiles) == 0 {
		return true
	}
	activeProfiles := collections.SliceToSet(this.activeProfiles())
	processor := regex.PatternProcessorOf(regex.NewPatternBuilder().Next("{word:\\w+}|{sign:\\W}").Build())
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
func (this *Environment) activeProfiles() []string {
	result := slices.Clone(this.profiles)
	for _, profile := range this.includedProfiles {
		if !slices.Contains(result, profile) {
			result = append(result, profile)
		}
	}
	return result
}

// first wins
func (this *Environment) propertySources() []PropertySource {
	return collections.ReverseSlice(this.sources)
}

// PROFILES_ACTIVE=dev,hsqldb
// PROFILES_INCLUDE=kubernetes
func (this *Environment) loadEnvironmentVariables() {
	environ := MapPropertySourceOf("Environment variables")
	pattern := regexp.MustCompile(regex.NewPatternBuilder().Next(`{key:[^=\s]+}={value:.*}`).Build())
	for _, keyValue := range os.Environ() {
		for _, m := range pattern.FindAllStringSubmatchIndex(keyValue, -1) {
			match := regex.MatchOf(pattern, keyValue, m)
			environ.SetProperty(match.NamedGroup("key").Value(), match.NamedGroup("value").Value())
		}
	}
	this.environPropertySource = environ
}

// --profiles.active=dev,hsqldb
// --profiles.include=kubernetes
func (this *Environment) loadApplicationParameters() {
	params := MapPropertySourceOf("Application parameters")
	pattern := regexp.MustCompile(regex.NewPatternBuilder().Next(`--?{key:[^=\s]+}\s*=?{value:.*}`).Build())
	for _, keyValue := range os.Args[1:] {
		for _, m := range pattern.FindAllStringSubmatchIndex(keyValue, -1) {
			match := regex.MatchOf(pattern, keyValue, m)
			params.SetProperty(match.NamedGroup("key").Value(), match.NamedGroup("value").Value())
		}
	}
	this.paramsPropertySource = params
}

// last wins
// application.yaml
// application-<profile>.yaml
func (this *Environment) loadApplicationConfiguration(bootstrapProfiles string) {
	this.profiles = splitProfiles(objects.FirstNonZero(bootstrapProfiles, this.paramsPropertySource.properties[profilesActiveProp], this.environPropertySource.properties[profilesActiveEnv]))
	this.includeProfiles(this.environPropertySource.properties[profilesIncludeEnv])
	this.includeProfiles(this.paramsPropertySource.properties[profilesIncludeProp])

	configName := objects.FirstNonZero(this.paramsPropertySource.properties["config.name"], this.environPropertySource.properties["CONFIG_NAME"], "application")
	defaultLocation := this.defaultConfigLocation()
	additionalLocation := objects.FirstNonZero(this.paramsPropertySource.properties["config.additional-location"], this.environPropertySource.properties["CONFIG_ADDITIONALLOCATION"])
	extendedDefaultLocation := lang.If(len(additionalLocation) == 0, defaultLocation, defaultLocation+","+additionalLocation)
	configLocation := objects.FirstNonZero(this.paramsPropertySource.properties["config.location"], this.environPropertySource.properties["CONFIG_LOCATION"])
	extendedConfigLocation := lang.If(len(additionalLocation) == 0, configLocation, additionalLocation+","+configLocation)
	resolvedConfigLocation := lang.If(len(configLocation) == 0, extendedDefaultLocation, extendedConfigLocation)

	this.discoverProfiles(resolvedConfigLocation, configName)
	this.sources = this.sources[:0]
	this.loadConfigurations(resolvedConfigLocation, configName)
	for _, source := range this.sources {
		slog.Info(fmt.Sprintf("Loaded configuration from %s", source.Name()))
	}
}

func (this *Environment) discoverProfiles(resolvedConfigLocation, configName string) {
	for {
		this.sources = this.sources[:0]
		profileCount := len(this.activeProfiles())
		this.loadConfigurations(resolvedConfigLocation, configName)
		if len(this.activeProfiles()) == profileCount {
			return
		}
	}
}

func (this *Environment) loadConfigurations(resolvedConfigLocation, configName string) {
	profiles := this.activeProfiles()
	for _, location := range strings.Split(resolvedConfigLocation, ",") {
		for i := 0; i <= len(profiles); i++ {
			profile := defaultProfile
			if i > 0 {
				profile = profiles[i-1]
			}
			for _, locationGroup := range strings.Split(location, ";") {
				this.loadConfiguration(locationGroup, configName, profile)
			}
		}
	}
}

func (this *Environment) defaultConfigLocation() string {
	if !this.isTest() {
		return "./,./config/"
	}
	root := this.moduleRoot()
	return filepath.ToSlash(root) + "/," + filepath.ToSlash(filepath.Join(root, "config")) + "/"
}

func (this *Environment) loadConfiguration(location, name, profile string) {
	location = filepath.ToSlash(location)
	var fantomExt string
	for _, m := range locationPattern.FindAllStringSubmatchIndex(location, -1) {
		match := regex.MatchOf(locationPattern, location, m)
		location = match.NamedGroup("location").Value()
		fantomExt = match.NamedGroup("fantomExt").Value()
	}

	if strings.HasSuffix(location, "/") {
		this.loadFile(files.RelativePath(location, lang.If(profile == defaultProfile, name+".yml", name+"-"+profile+".yml")), fantomExt)
		this.loadFile(files.RelativePath(location, lang.If(profile == defaultProfile, name+".yaml", name+"-"+profile+".yaml")), fantomExt)
		this.loadFile(files.RelativePath(location, lang.If(profile == defaultProfile, name+".properties", name+"-"+profile+".properties")), fantomExt)
	} else if len(fantomExt) > 0 {
		this.loadFile(lang.If(profile == defaultProfile, location, location+"-"+profile), fantomExt)
	} else {
		ext := filepath.Ext(location)
		this.loadFile(lang.If(profile == defaultProfile, location, location[:len(location)-len(ext)]+"-"+profile+ext), fantomExt)
	}
}

func (this *Environment) loadFile(path, fantomExt string) {
	if !files.Exists(path) {
		return
	}
	var result PropertySource
	ext := objects.FirstNonZero(fantomExt, filepath.Ext(path))
	lang.Assert(len(ext) != 0, "Cannot load from location %s. If location supposed to be a directory use '/' at the end. Otherwise provide extension hint in square brackets like [.properties] to derive property source type", path)
	file := optional.OfCommaErr(os.Open(path)).OrElsePanic("Cannot open file %s", path)
	defer file.Close()
	content := string(optional.OfCommaErr(io.ReadAll(file)).OrElsePanic("Cannot read from %s", path))
	switch ext {
	case ".properties":
		result = NewPropertiesPropertySource(path, content)
	case ".yaml", ".yml":
		result = NewYamlPropertySource(path, content)
	default:
		panic(err.NewRuntimeException(fmt.Sprintf("Cannot load from %s as %s file type is not supported. Use extension hint in square brackets like .env[.properties] to derive property source type", path, ext)))
	}
	this.sources = append(this.sources, result)
	if result.HasProperty(profilesActiveProp) && len(this.profiles) == 0 {
		this.profiles = splitProfiles(result.Property(profilesActiveProp))
	}
	if result.HasProperty(profilesIncludeProp) {
		this.includeProfiles(result.Property(profilesIncludeProp))
	}
	if result.HasProperty("config.import") {
		for _, location := range strings.Split(result.Property("config.import"), ",") {
			this.loadImport(path, location)
		}
	}
}

func (this *Environment) loadImport(path, location string) {
	var fantomExt string
	for _, m := range locationPattern.FindAllStringSubmatchIndex(location, -1) {
		match := regex.MatchOf(locationPattern, location, m)
		location = match.NamedGroup("location").Value()
		fantomExt = match.NamedGroup("fantomExt").Value()
	}
	location = filepath.ToSlash(location)
	lang.Assert(!strings.HasSuffix(location, "/"), "Cannot load from location %s defined in %s. Directory import is not supported", location, path)
	this.loadFile(files.RelativePath(path, location), fantomExt)
}

func (this *Environment) includeProfiles(profiles string) {
	for _, profile := range splitProfiles(profiles) {
		if !slices.Contains(this.includedProfiles, profile) {
			this.includedProfiles = append(this.includedProfiles, profile)
		}
	}
}

func (this *Environment) envVarCanonicalForm(key string) string {
	return strings.ToUpper(str.ReplaceChars(key, envVarCanonicalFormTranslationRule))
}

func (this *Environment) isTest() bool {
	return stream.From(os.Args[1:]).
		Filter(func(s string) bool { return strings.HasPrefix(s, "-test.timeout=") }).
		FindFirst().Present()
}

func (this *Environment) moduleRoot() string {
	dir := optional.OfCommaErr(os.Getwd()).OrElsePanic("Cannot get working directory")
	for {
		if files.Exists(filepath.Join(dir, "go.mod")) {
			return dir
		}
		parent := filepath.Dir(dir)
		lang.Assert(parent != dir, "Cannot find go.mod from %q", dir)
		dir = parent
	}
}

func (this *Environment) addPropertySource(source PropertySource) {
	this.sources = append(this.sources, source)
	slog.Debug(fmt.Sprintf("Added property source %T", source))
}

func (this *Environment) setContextVariable(key string, value any) {
	this.exprProcessor.Define(key, value)
	this.strictExprProcessor.Define(key, value)
}
