package io

import (
	"io"
	"net/url"
	"os"
	"time"

	"github.com/go-external-config/v1/util"
)

type FileResource struct {
	url *url.URL
}

func NewFileResource(url *url.URL) *FileResource {
	return &FileResource{
		url: url}
}

func (r *FileResource) Reader() io.Reader {
	return util.OptionalOfCommaErr(os.Open(r.url.Path)).OrElsePanic("Cannot open file %s", r.url.Path)
}

func (r *FileResource) URL() *url.URL {
	return r.url
}

func (r *FileResource) CreateRelative(location string) Resource {
	relativeLocation := util.OptionalOfCommaErr(url.PathUnescape(r.url.JoinPath(location).String())).OrElsePanic("Cannot create relative location %s + %s", r.url.Path, location)
	relativeUrl := util.OptionalOfCommaErr(url.Parse(relativeLocation)).OrElsePanic("Cannot parse location %s", relativeLocation)
	return NewFileResource(relativeUrl)
}

func (r *FileResource) Exists() bool {
	return util.OptionalOfCommaErr(os.Stat(r.url.Path)).OrElse(nil) != nil
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

func (r *FileResource) String() string {
	return util.OptionalOfCommaErr(url.PathUnescape(r.url.String())).OrElsePanic("Cannot unescape %s", r.url)
}
