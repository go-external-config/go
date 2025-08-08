package io

import (
	"fmt"
	"net/url"

	"github.com/go-external-config/v1/lang"
	"github.com/go-external-config/v1/util"
)

type ProtocolResolver interface {
	Resolve(location string) Resource
}

type ResourceLoader struct {
	protocolResolvers []ProtocolResolver
}

func NewResourceLoader() *ResourceLoader {
	resourceLoader := ResourceLoader{
		protocolResolvers: make([]ProtocolResolver, 1)}
	resourceLoader.protocolResolvers = append(resourceLoader.protocolResolvers, &resourceLoader)
	return &resourceLoader
}

func (l *ResourceLoader) Resource(location string) Resource {
	for _, protocolResolver := range l.protocolResolvers {
		resource := protocolResolver.Resolve(location)
		if resource != nil {
			return resource
		}
	}
	panic("Cannot resolve resource " + location)
}

func (l *ResourceLoader) Resolve(location string) Resource {
	url := util.OptionalOfCommaErr(url.Parse(location)).OrElsePanic("Cannot parse location %s", location)
	schema := lang.FirstNonEmpty(url.Scheme, "file")
	switch schema {
	case "file":
		return NewFileResource(url)
	case "embed":
		return NewEmbedResource(url)
	}
	panic(fmt.Sprintf("Cannot resolve resource of schema '%s' in %s", schema, location))
}

func (l *ResourceLoader) AddProtocalResolver(protocolResolver ProtocolResolver) {
	l.protocolResolvers = append(l.protocolResolvers, protocolResolver)
}
