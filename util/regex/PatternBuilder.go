package regex

import (
	"fmt"
	"strings"

	"github.com/go-external-config/go/lang"
)

type PatternBuilder struct {
	PatternProcessor
	buf strings.Builder
}

func NewPatternBuilder(mod ...string) *PatternBuilder {
	this := &PatternBuilder{}
	this.PatternProcessor = *PatternProcessorOf(`(?P<space>[ \t]+)|(?P<preffix>[^\\])?\{(?P<param>.*?)(?P<suffix>[^\\]|\{\d+,?\d*[^\d])?\}`)
	this.OverrideResolve(this.resolve)
	mod = append(mod, "(?ms)")
	this.buf.WriteString(mod[0])
	return this
}

func (this *PatternBuilder) Start(pattern string) *PatternBuilder {
	this.buf.WriteString(`^[ \t]*?`)
	this.buf.WriteString(this.ProcessRecursive(pattern, false).(string))
	return this
}

func (this *PatternBuilder) End() *PatternBuilder {
	this.buf.WriteString(`[ \t]*\r?(\n|$)`)
	return this
}

func (this *PatternBuilder) Next(pattern string) *PatternBuilder {
	this.buf.WriteString(this.ProcessRecursive(pattern, false).(string))
	return this
}

func (this *PatternBuilder) Opt(pattern string) *PatternBuilder {
	this.buf.WriteString(`(?:`)
	this.buf.WriteString(this.ProcessRecursive(pattern, false).(string))
	this.buf.WriteString(`)?`)
	return this
}

func (this *PatternBuilder) Build() string {
	return this.buf.String()
}

func (this *PatternBuilder) resolve(match *Match, super func(*Match) any) any {
	space := match.NamedGroup("space")
	if space.Present() {
		return `[ \t]+`
	}

	preffix := match.NamedGroup("preffix")
	param := match.NamedGroup("param")
	suffix := match.NamedGroup("suffix")
	paramStr := param.Value() + suffix.OrElse("")
	var paramRegex string
	regexIdx := strings.Index(paramStr, ":")
	if regexIdx > 0 {
		paramRegex = paramStr[regexIdx+1:]
		paramStr = paramStr[:regexIdx]
	}

	return fmt.Sprintf(`%s(?P<%s>%s)`, preffix.OrElse(""), paramStr, lang.If(paramRegex == "", ".+?", paramRegex))
}
