package env

import (
	"os"
	"regexp"

	"github.com/go-jang/go/util/regex"
)

// --profiles.active=dev,hsqldb
// --profiles.include=kubernetes
type ParamsPropertySource struct {
	*MapPropertySource
}

func NewParamsPropertySource() *ParamsPropertySource {
	this := &ParamsPropertySource{
		MapPropertySource: MapPropertySourceOf("Application parameters"),
	}

	pattern := regexp.MustCompile(regex.NewPatternBuilder().Next(`--?{key:[^=\s]+}\s*=?{value:.*}`).Build())

	for _, keyValue := range os.Args[1:] {
		for _, m := range pattern.FindAllStringSubmatchIndex(keyValue, -1) {
			match := regex.MatchOf(pattern, keyValue, m)
			this.SetProperty(match.NamedGroup("key").Value(), match.NamedGroup("value").Value())
		}
	}
	return this
}
