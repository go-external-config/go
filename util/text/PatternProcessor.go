package text

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-external-config/go/lang"
	"github.com/go-external-config/go/util/regex"
)

type PatternProcessor struct {
	regexp  *regexp.Regexp
	resolve func(*regex.Match) any
}

func PatternProcessorOf(pattern string) *PatternProcessor {
	return &PatternProcessor{
		regexp:  regexp.MustCompile(lang.If(strings.HasPrefix(pattern, "(?"), pattern, "(?ms)"+pattern)),
		resolve: func(*regex.Match) any { panic("Not implemented") }}
}

func (p *PatternProcessor) Process(str string) any {
	return p.ProcessRecursive(str, true)
}

func (p *PatternProcessor) ProcessRecursive(str string, recursive bool) any {
	if str == "" {
		return str
	}
	before, resolved := str, str
	for {
		var sb strings.Builder
		before = resolved
		matched := p.regexp.FindAllStringSubmatchIndex(resolved, -1)
		lastMatch := []int{0, 0}
		for _, match := range matched {
			regexMatch := regex.MatchOf(p.regexp, resolved, match)
			// full match may be resolved to non-string value in subclasses
			if len(matched) == 1 && len(regexMatch.Expr()) == len(resolved) {
				candidate := p.resolve(regexMatch)
				switch v := candidate.(type) {
				case string:
					sb.WriteString(v)
				default:
					return v
				}
			} else {
				sb.WriteString(resolved[lastMatch[1]:match[0]])
				sb.WriteString(fmt.Sprint(p.resolve(regexMatch)))
			}
			lastMatch = match
		}
		sb.WriteString(resolved[lastMatch[1]:])
		resolved = sb.String()
		if !recursive || before == resolved {
			break
		}
	}
	return resolved
}

func (p *PatternProcessor) OverrideResolve(f func(*regex.Match,
	func(*regex.Match) any) any) {

	super := p.resolve
	p.resolve = func(rm *regex.Match) any { return f(rm, super) }
}
