package text

import (
	"regexp"
	"strings"

	"github.com/madamovych/go/lang"
)

type PatternProcessor struct {
	regexp  *regexp.Regexp
	resolve func(matcher string) string
}

func NewtPatternProcessor(pattern string) *PatternProcessor {
	return &PatternProcessor{
		regexp:  regexp.MustCompile(lang.If(strings.HasPrefix(pattern, "(?"), pattern, "(?ms)"+pattern)),
		resolve: func(matcher string) string { panic("Not implemented") }}
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
		matched := p.regexp.FindAllStringIndex(resolved, -1)
		lastMatch := []int{0, 0}
		for idx, match := range matched {
			sb.WriteString(resolved[lastMatch[1]:match[0]])
			sb.WriteString(p.resolve(resolved[match[0]:match[1]]))
			if idx == len(matched)-1 {
				sb.WriteString(resolved[match[1]:])
			}
			lastMatch = match
		}
		resolved = sb.String()
		if !recursive || before == resolved {
			break
		}
	}
	return resolved
}
