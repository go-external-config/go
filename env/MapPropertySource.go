package env

import (
	"github.com/madamovych/go/util"
	"github.com/madamovych/go/util/text"
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
	return util.OptionalOfEntry(s.properties, key).Present()
}

func (s *MapPropertySource) Property(key string) string {
	return util.OptionalOfEntry(s.properties, key).
		OrElsePanic("%v has no %v", s.name, key)
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

func (s *MapPropertySource) ResolvePlaceholders() {
	processor := text.ExprProcessorOf(false)
	for todo := true; todo; {
		todo = false
		for key, value := range s.properties {
			processor.Define(key, value)
			resolved := processor.Process(value)
			if resolved != value {
				todo = true
				s.properties[key] = resolved
				processor.Define(key, resolved)
				// fmt.Printf("PropertySource[%v]: %v: %v -> %v\n", s.name, key, value, resolved)
			}
		}
	}
	processor.SetStrict(true)
	for _, value := range s.properties {
		processor.Process(value)
	}
}
