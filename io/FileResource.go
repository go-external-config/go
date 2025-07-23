package io

import (
	"io"
	"net/url"
	"os"
	"time"

	"github.com/madamovych/go/lang"
)

type FileResource struct {
	url *url.URL
}

func NewFileResource(url *url.URL) *FileResource {
	return &FileResource{
		url: url}
}

func (r *FileResource) Reader() io.Reader {
	return lang.OptionalOfCommaErr(os.Open(r.url.Path)).OrElsePanic("Cannot open file %s", r.url.Path)
}

func (r *FileResource) URL() *url.URL {
	return r.url
}

func (r *FileResource) CreateRelative(location string) Resource {
	relativeUrl := r.URL().JoinPath(location)
	return NewFileResource(relativeUrl)
}

func (r *FileResource) Exists() bool {
	return lang.OptionalOfCommaErr(os.Stat(r.url.Path)).OrElse(nil) != nil
}

func (r *FileResource) Name() string {
	panic("Not implemented")
}

func (r *FileResource) IsDir() bool {
	panic("Not implemented")
}

func (r *FileResource) Size() int64 {
	panic("Not implemented")
}

func (r *FileResource) ModTime() time.Time {
	panic("Not implemented")
}
