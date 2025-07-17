package text

import (
	"regexp"
	"strings"

	"github.com/madamovych/go/lang"
)

type PatternProcessor struct {
	regexp  *regexp.Regexp
	resolve func(*lang.RegexpMatch) string
}

func NewPatternProcessor(pattern string) *PatternProcessor {
	return &PatternProcessor{
		regexp:  regexp.MustCompile(lang.If(strings.HasPrefix(pattern, "(?"), pattern, "(?ms)"+pattern)),
		resolve: func(*lang.RegexpMatch) string { panic("Not implemented") }}
}

func (p *PatternProcessor) Process(str string) string {
	return p.ProcessRecursive(str, true)
}

func (p *PatternProcessor) ProcessRecursive(str string, recursive bool) string {
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
			sb.WriteString(resolved[lastMatch[1]:match[0]])
			sb.WriteString(p.resolve(lang.NewRegexpMatch(p.regexp, resolved, match)))
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

func (p *PatternProcessor) SetResolve(f func(*lang.RegexpMatch) string) {
	p.resolve = f
}
