package env

import (
	"fmt"

	"github.com/go-errr/go/err"
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

func (this *YamlPropertySource) propertiesFromYaml(yamlStr string) map[string]string {
	var parsedYaml any
	e := yaml.Unmarshal([]byte(yamlStr), &parsedYaml)
	if e != nil {
		panic(err.NewRuntimeException(fmt.Sprintf("Error unmarshaling yaml: %v", e)))
	}
	properties := make(map[string]string)
	this.flattenYaml(parsedYaml, "", properties)
	return properties
}

func (this *YamlPropertySource) flattenYaml(data any, prefix string, result map[string]string) {
	switch v := data.(type) {
	case map[string]any:
		for key, value := range v {
			newPrefix := key
			if prefix != "" {
				newPrefix = prefix + "." + key
			}
			this.flattenYaml(value, newPrefix, result)
		}
	case []any:
		for i, value := range v {
			newPrefix := fmt.Sprintf("%s[%d]", prefix, i)
			this.flattenYaml(value, newPrefix, result)
		}
	default:
		result[prefix] = fmt.Sprint(v)
	}
}
