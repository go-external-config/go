package env

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type YamlPropertySource struct {
	MapPropertySource
}

func NewYamlPropertySource(name, yaml string) *YamlPropertySource {
	yamlPropertySource := YamlPropertySource{
		MapPropertySource: *MapPropertySourceOf(name)}
	yamlPropertySource.SetProperties(yamlPropertySource.propertiesFromYaml(yaml))
	return &yamlPropertySource
}

func (s *YamlPropertySource) propertiesFromYaml(yamlStr string) map[string]string {
	var parsedYaml any
	err := yaml.Unmarshal([]byte(yamlStr), &parsedYaml)
	if err != nil {
		panic(fmt.Sprintf("Error unmarshaling yaml: %v\n", err))
	}
	properties := make(map[string]string)
	s.flattenYaml(parsedYaml, "", properties)
	return properties
}

func (s *YamlPropertySource) flattenYaml(data any, prefix string, result map[string]string) {
	switch v := data.(type) {
	case map[string]any:
		for key, value := range v {
			newPrefix := key
			if prefix != "" {
				newPrefix = prefix + "." + key
			}
			s.flattenYaml(value, newPrefix, result)
		}
	case []any:
		for i, value := range v {
			newPrefix := fmt.Sprintf("%s[%d]", prefix, i)
			s.flattenYaml(value, newPrefix, result)
		}
	default:
		result[prefix] = fmt.Sprint(v)
	}
}
