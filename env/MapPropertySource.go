package env

import (
	"github.com/go-external-config/go/lang"
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

func (this *MapPropertySource) Name() string {
	return this.name
}

func (this *MapPropertySource) HasProperty(key string) bool {
	_, ok := this.properties[key]
	return ok
}

func (this *MapPropertySource) Property(key string) string {
	value, ok := this.properties[key]
	lang.AssertState(ok, "%v has no %v", this.name, key)
	return value
}

func (this *MapPropertySource) SetProperty(key string, value string) {
	this.properties[key] = value
}

func (this *MapPropertySource) Properties() map[string]string {
	return this.properties
}

func (this *MapPropertySource) SetProperties(properties map[string]string) {
	this.properties = properties
}

func (this *MapPropertySource) ContainsProperty(key string) bool {
	_, ok := this.properties[key]
	return ok
}
