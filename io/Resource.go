package io

import (
	"io"
	"net/url"
	"time"
)

type Resource interface {
	Reader() io.Reader
	URL() *url.URL
	CreateRelative(location string) Resource
	Exists() bool
	Name() string
	IsDir() bool
	Size() int64
	ModTime() time.Time
	String() string
}
