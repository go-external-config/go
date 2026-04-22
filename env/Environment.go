package env

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/go-errr/go/err"
	"github.com/go-external-config/go/lang"
	"github.com/go-external-config/go/util/collection"
	"github.com/go-external-config/go/util/concurrent"
	"github.com/go-external-config/go/util/files"
	"github.com/go-external-config/go/util/optional"
	"github.com/go-external-config/go/util/regex"
	"github.com/go-external-config/go/util/str"
)

var locationPattern = regexp.MustCompile(regex.NewPatternBuilder().Next(`{location:.+}\[{fantomExt:\.[\w]+}\]`).Build())
var envVarCanonicalFormTranslationRule = map[rune]rune{
	'.': '_',
	'[': '_',
	']': '_',
	'-': 0, // delete
}

var environment *Environment
var environmentMu sync.Mutex

type Environment struct {
	activeProfiles        []string
	paramsPropertySource  *MapPropertySource
	environPropertySource *MapPropertySource
	propertySources       []PropertySource
	exprProcessor         *ExprProcessor
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
		activeProfiles:  []string{"default"},
		propertySources: make([]PropertySource, 0),
		exprProcessor:   ExprProcessorOf(true)}

	environment.loadEnvironmentVariables()
	environment.loadApplicationParameters()
	environment.loadApplicationConfiguration(activeProfiles)
	environment.WithPropertySource(NewRandomValuePropertySource())
	environment.WithPropertySource(NewBase64PropertySource(&environment))
	return &environment
}

func (this *Environment) Property(key string) string {
	return fmt.Sprint(this.ResolveRequiredPlaceholders(this.lookupRawProperty(key).
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
		for i := len(this.propertySources) - 1; i >= 0; i-- {
			if this.propertySources[i].HasProperty(key) {
				return optional.OfValue(this.propertySources[i].Property(key))
			}
		}
	}
	return optional.OfEmpty[string]()
}

func (this *Environment) ResolveRequiredPlaceholders(expression string) any {
	return this.exprProcessor.Process(expression)
}

// Determine whether one or more of the given profiles is active.
//
// If a profile begins with '!' the logic is inverted, meaning this method will return true if the given profile is not active.
// For example, env.MatchesProfiles("p1", "!p2") will return true if profile 'p1' is active or 'p2' is not active.
// A compound expression allows for more complicated profile logic to be expressed, for example "production & cloud".
func (this *Environment) MatchesProfiles(profiles ...string) bool {
	if len(profiles) == 0 {
		return true
	}
	activeProfiles := collection.SliceToSet(this.activeProfiles)
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
func (this *Environment) ActiveProfiles() []string {
	return this.activeProfiles
}

// first wins
func (this *Environment) PropertySources() []PropertySource {
	return collection.ReverseSlice(this.propertySources)
}

// PROFILES_ACTIVE=dev,hsqldb
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
	activeProfiles := lang.FirstNonEmpty(bootstrapProfiles, this.paramsPropertySource.properties["profiles.active"], this.environPropertySource.properties["PROFILES_ACTIVE"])
	this.activeProfiles = lang.If(len(activeProfiles) == 0, this.activeProfiles, append(this.activeProfiles, strings.Split(activeProfiles, ",")...))
	configName := lang.FirstNonEmpty(this.paramsPropertySource.properties["config.name"], this.environPropertySource.properties["CONFIG_NAME"], "application")
	defaultLocation := "./,./config/"
	additionalLocation := lang.FirstNonEmpty(this.paramsPropertySource.properties["config.additional-location"], this.environPropertySource.properties["CONFIG_ADDITIONALLOCATION"])
	extendedDefaultLocation := lang.If(len(additionalLocation) == 0, defaultLocation, defaultLocation+","+additionalLocation)
	configLocation := lang.FirstNonEmpty(this.paramsPropertySource.properties["config.location"], this.environPropertySource.properties["CONFIG_LOCATION"])
	extendedConfigLocation := lang.If(len(additionalLocation) == 0, configLocation, additionalLocation+","+configLocation)
	resolvedConfigLocation := lang.If(len(configLocation) == 0, extendedDefaultLocation, extendedConfigLocation)

	for _, location := range strings.Split(resolvedConfigLocation, ",") {
		for i := 0; i < len(this.activeProfiles); i++ {
			for _, locationGroup := range strings.Split(location, ";") {
				this.loadConfiguration(locationGroup, configName, this.activeProfiles[i])
			}
		}
	}
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
		this.loadFile(files.RelativePath(location, lang.If(profile == "default", name+".yml", name+"-"+profile+".yml")), fantomExt)
		this.loadFile(files.RelativePath(location, lang.If(profile == "default", name+".yaml", name+"-"+profile+".yaml")), fantomExt)
		this.loadFile(files.RelativePath(location, lang.If(profile == "default", name+".properties", name+"-"+profile+".properties")), fantomExt)
	} else if len(fantomExt) > 0 {
		this.loadFile(lang.If(profile == "default", location, location+"-"+profile), fantomExt)
	} else {
		ext := filepath.Ext(location)
		this.loadFile(lang.If(profile == "default", location, location[:len(location)-len(ext)]+"-"+profile+ext), fantomExt)
	}
}

func (this *Environment) loadFile(path, fantomExt string) {
	if !files.Exists(path) {
		return
	}
	var result PropertySource
	fmt.Printf("loading properties from %s\n", path)
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
		panic(err.NewRuntimeException(fmt.Sprintf("Cannot load from %s as %s file type is not supported. Use extension hint in square brackets like .env[.properties] to derive property source type", path, ext)))
	}
	this.propertySources = append(this.propertySources, result)
	if result.HasProperty("profiles.active") && len(this.activeProfiles) == 1 && this.activeProfiles[0] == "default" {
		this.activeProfiles = append(this.activeProfiles, strings.Split(result.Property("profiles.active"), ",")...)
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
	lang.AssertState(!strings.HasSuffix(location, "/"), "Cannot load from location %s defined in %s. Directory import is not supported", location, path)
	this.loadFile(files.RelativePath(path, location), fantomExt)
}

func (this *Environment) envVarCanonicalForm(key string) string {
	return strings.ToUpper(str.ReplaceChars(key, envVarCanonicalFormTranslationRule))
}

// Add custom property source to implement additional logic for properties processing, like property=base64:dGVzdAo=.
// See Base64PropertySource (available by default) and RsaPropertySource
//
//	var _ = env.Instance().WithPropertySource(env.NewRsaPropertySource())
func (this *Environment) WithPropertySource(source PropertySource) *Environment {
	this.propertySources = append(this.propertySources, source)
	return this
}

// Add custom context variables to be evaluated.
// See env.ExprProcessor for expressions and variables available by default.
//
//	var _ = env.Instance().WithContextVariable("runtime", map[string]any{
//		"NumCPU": runtime.NumCPU(),
//	})
func (this *Environment) WithContextVariable(key string, value any) *Environment {
	this.exprProcessor.Define(key, value)
	return this
}
