package env

import (
	"encoding/base64"
	"strings"

	"github.com/go-external-config/go/util/optional"
)

// Custom property source as an additional logic for properties processing, like property=base64:dGVzdAo=
type Base64PropertySource struct {
	environment *Environment
}

func NewBase64PropertySource(environment *Environment) *Base64PropertySource {
	return &Base64PropertySource{
		environment: environment}
}

func (s *Base64PropertySource) Name() string {
	return "Base64PropertySource"
}

func (s *Base64PropertySource) HasProperty(key string) bool {
	for _, source := range environment.PropertySources() {
		if source.Properties() != nil && source.HasProperty(key) {
			return strings.HasPrefix(source.Property(key), "base64:")
		}
	}
	return false
}

func (s *Base64PropertySource) Property(key string) string {
	for _, source := range environment.PropertySources() {
		if source.Properties() != nil && source.HasProperty(key) {
			value := source.Property(key)[7:]
			return strings.TrimRight(string(optional.OfCommaErr(base64.StdEncoding.DecodeString(value)).
				OrElsePanic("Cannot decode %s=%s", key, value)), "\n\r")
		}
	}
	panic("No value present for " + key)
}

func (s *Base64PropertySource) Properties() map[string]string {
	return nil
}
