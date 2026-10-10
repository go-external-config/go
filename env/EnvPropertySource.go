package env

import (
	"os"
	"regexp"
	"strings"

	"github.com/go-external-config/go/str"
	"github.com/go-jang/go/util/regex"
)

var envVarCanonicalFormTranslationRule = map[rune]rune{
	'.': '_',
	'[': '_',
	']': '_',
	'-': 0,
}

// PROFILES_ACTIVE=dev,hsqldb
// PROFILES_INCLUDE=kubernetes
type EnvPropertySource struct {
	*MapPropertySource
}

func NewEnvPropertySource() *EnvPropertySource {
	this := &EnvPropertySource{
		MapPropertySource: MapPropertySourceOf("Environment variables"),
	}

	pattern := regexp.MustCompile(regex.NewPatternBuilder().Next(`{key:[^=\s]+}={value:.*}`).Build())

	for _, keyValue := range os.Environ() {
		for _, m := range pattern.FindAllStringSubmatchIndex(keyValue, -1) {
			match := regex.MatchOf(pattern, keyValue, m)
			this.SetProperty(match.NamedGroup("key").Value(), match.NamedGroup("value").Value())
		}
	}
	return this
}

func (this *EnvPropertySource) HasProperty(key string) bool {
	return this.MapPropertySource.HasProperty(key) ||
		this.MapPropertySource.HasProperty(this.canonicalForm(key))
}

func (this *EnvPropertySource) Property(key string) string {
	if this.MapPropertySource.HasProperty(key) {
		return this.MapPropertySource.Property(key)
	}
	canonical := this.canonicalForm(key)
	if this.MapPropertySource.HasProperty(canonical) {
		return this.MapPropertySource.Property(canonical)
	}
	return this.MapPropertySource.Property(key)
}

func (this *EnvPropertySource) canonicalForm(key string) string {
	return strings.ToUpper(str.ReplaceChars(key, envVarCanonicalFormTranslationRule))
}
