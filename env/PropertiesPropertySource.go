package env

import "github.com/magiconair/properties"

type PropertiesPropertySource struct {
	MapPropertySource
}

func NewPropertiesPropertySource(name, content string) *PropertiesPropertySource {
	propertiesPropertySource := PropertiesPropertySource{
		MapPropertySource: *MapPropertySourceOf(name)}
	propertiesPropertySource.SetProperties(propertiesPropertySource.propertiesFrom(content))
	return &propertiesPropertySource
}

func (s *PropertiesPropertySource) propertiesFrom(content string) map[string]any {
	result := make(map[string]any)
	for key, value := range properties.MustLoadString(content).Map() {
		result[key] = value
	}
	return result
}
