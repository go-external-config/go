package io

import (
	"net/url"

	"github.com/madamovych/go/lang"
	"github.com/madamovych/go/util"
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
	return nil
}

func (l *ResourceLoader) AddProtocalResolver(protocolResolver ProtocolResolver) {
	l.protocolResolvers = append(l.protocolResolvers, protocolResolver)
}
