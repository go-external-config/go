package env

import (
	"github.com/madamovych/go/lang"
	"github.com/madamovych/go/text"
)

type MapPropertySource struct {
	name       string
	properties map[string]any
}

func MapPropertySourceOf(name string) *MapPropertySource {
	return &MapPropertySource{
		name:       name,
		properties: make(map[string]any)}
}

func MapPropertySourceOfMap(name string, source map[string]any) *MapPropertySource {
	return &MapPropertySource{
		name:       name,
		properties: source}
}

func (s *MapPropertySource) Name() string {
	return s.name
}

func (s *MapPropertySource) Property(key string) any {
	return lang.OptionalOfEntry(s.properties, key).
		OrElsePanic("%v has no %v", s.name, key)
}

func (s *MapPropertySource) SetProperty(key string, value any) {
	s.properties[key] = value
}

func (s *MapPropertySource) Properties() map[string]any {
	return s.properties
}

func (s *MapPropertySource) SetProperties(properties map[string]any) {
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
			switch strValue := value.(type) {
			case string:
				resolved := processor.Process(strValue)
				if resolved != strValue {
					todo = true
					s.properties[key] = resolved
					processor.Define(key, resolved)
					// fmt.Printf("PropertySource[%v]: %v: %v -> %v\n", s.name, key, value, resolved)
				}
			}
		}
	}
	processor.SetStrict(true)
	for _, value := range s.properties {
		switch strValue := value.(type) {
		case string:
			processor.Process(strValue)
		}
	}
}
