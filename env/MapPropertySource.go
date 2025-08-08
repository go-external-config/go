package env

import (
	"github.com/go-external-config/v1/lang"
)

type MapPropertySource struct {
	name       string
	properties map[string]string
}

func MapPropertySourceOf(name string) *MapPropertySource {
	return &MapPropertySource{
		name:       name,
		properties: make(map[string]string)}
}

func MapPropertySourceOfMap(name string, source map[string]string) *MapPropertySource {
	return &MapPropertySource{
		name:       name,
		properties: source}
}

func (s *MapPropertySource) Name() string {
	return s.name
}

func (s *MapPropertySource) HasProperty(key string) bool {
	_, ok := s.properties[key]
	return ok
}

func (s *MapPropertySource) Property(key string) string {
	value, ok := s.properties[key]
	lang.AssertState(ok, "%v has no %v", s.name, key)
	return value
}

func (s *MapPropertySource) SetProperty(key string, value string) {
	s.properties[key] = value
}

func (s *MapPropertySource) Properties() map[string]string {
	return s.properties
}

func (s *MapPropertySource) SetProperties(properties map[string]string) {
	s.properties = properties
}

func (s *MapPropertySource) ContainsProperty(key string) bool {
	_, ok := s.properties[key]
	return ok
}
