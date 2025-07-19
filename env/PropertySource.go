package env

import (
	"github.com/madamovych/go/lang"
	"github.com/madamovych/go/text"
)

type PropertySource struct {
	name   string
	source map[string]any
}

func PropertySourceOf(name string) *PropertySource {
	return &PropertySource{
		name:   name,
		source: make(map[string]any)}
}

func PropertySourceOfMap(name string, source map[string]any) *PropertySource {
	return &PropertySource{
		name:   name,
		source: source}
}

func (s *PropertySource) Name() string {
	return s.name
}

func (s *PropertySource) Property(key string) any {
	return lang.OptionalOfEntry(s.source, key).
		OrElsePanic("%v has no %v", s.name, key)
}

func (s *PropertySource) SetProperty(key string, value any) {
	s.source[key] = value
}

func (s *PropertySource) ContainsProperty(key string) bool {
	_, ok := s.source[key]
	return ok
}

func (s *PropertySource) ResolvePlaceholders() {
	processor := text.ExprProcessorOf(false)
	for todo := true; todo; {
		todo = false
		for key, value := range s.source {
			processor.Define(key, value)
			switch strValue := value.(type) {
			case string:
				resolved := processor.Process(strValue)
				if resolved != strValue {
					todo = true
					s.source[key] = resolved
					processor.Define(key, resolved)
					// fmt.Printf("PropertySource[%v]: %v: %v -> %v\n", s.name, key, value, resolved)
				}
			}
		}
	}
	processor.SetStrict(true)
	for _, value := range s.source {
		switch strValue := value.(type) {
		case string:
			processor.Process(strValue)
		}
	}
}
