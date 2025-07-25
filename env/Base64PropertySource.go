package env

import (
	"encoding/base64"
	"strings"

	"github.com/go-external-config/go/util"
)

// custom property source as additional logic for properties processing, like property=Base64:dGVzdAo=
type Base64PropertySource struct {
	environment PropertySource
}

func NewBase64PropertySource() *Base64PropertySource {
	return &Base64PropertySource{
		environment: GetEnvironment()}
}

func (s *Base64PropertySource) Name() string {
	return "Base64PropertySource"
}

func (s *Base64PropertySource) HasProperty(key string) bool {
	for _, source := range environment.PropertySources() {
		if source.Properties() != nil && source.HasProperty(key) {
			return strings.HasPrefix(source.Property(key), "Base64:")
		}
	}
	return false
}

func (s *Base64PropertySource) Property(key string) string {
	for _, source := range environment.PropertySources() {
		if source.Properties() != nil && source.HasProperty(key) {
			value := source.Property(key)
			return strings.TrimRight(string(util.OptionalOfCommaErr(base64.StdEncoding.DecodeString(value[7:])).
				OrElsePanic("Cannot decode %s: %s", key, value)), "\n\r")
		}
	}
	panic("No value present for " + key)
}

func (s *Base64PropertySource) Properties() map[string]string {
	return nil
}
