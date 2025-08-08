package io

import (
	"bytes"
	"io"
	"net/url"
	"time"

	"github.com/go-external-config/v1/util"
)

type EmbedResource struct {
	url *url.URL
}

func NewEmbedResource(url *url.URL) *EmbedResource {
	return &EmbedResource{
		url: url}
}

func (r *EmbedResource) Reader() io.Reader {
	return bytes.NewReader([]byte("Embeded resource reader is not implemented"))
}

func (r *EmbedResource) URL() *url.URL {
	return r.url
}

func (r *EmbedResource) CreateRelative(location string) Resource {
	return nil
}

func (r *EmbedResource) Exists() bool {
	panic("Not implemented")
}

func (r *EmbedResource) Name() string {
	panic("Not implemented")
}

func (r *EmbedResource) IsDir() bool {
	panic("Not implemented")
}

func (r *EmbedResource) Size() int64 {
	panic("Not implemented")
}

func (r *EmbedResource) ModTime() time.Time {
	panic("Not implemented")
}

func (r *EmbedResource) String() string {
	return util.OptionalOfCommaErr(url.PathUnescape(r.url.String())).OrElsePanic("Cannot unescape %s", r.url)
}
