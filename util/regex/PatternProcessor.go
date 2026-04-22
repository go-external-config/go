package regex

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-errr/go/err"
	"github.com/go-external-config/go/lang"
)

type PatternProcessor struct {
	regexp  *regexp.Regexp
	resolve func(*Match) any
}

func PatternProcessorOf(pattern string) *PatternProcessor {
	return &PatternProcessor{
		regexp:  regexp.MustCompile(lang.If(strings.HasPrefix(pattern, "(?"), pattern, "(?ms)"+pattern)),
		resolve: func(*Match) any { panic(err.NewUnsupportedOperationException("Not implemented")) }}
}

func (this *PatternProcessor) Process(str string) any {
	return this.ProcessRecursive(str, true)
}

func (this *PatternProcessor) ProcessRecursive(str string, recursive bool) any {
	if str == "" {
		return str
	}
	//nolint
	before, resolved := str, str
	for {
		var sb strings.Builder
		before = resolved
		matched := this.regexp.FindAllStringSubmatchIndex(resolved, -1)
		lastMatch := []int{0, 0}
		for _, match := range matched {
			regexMatch := MatchOf(this.regexp, resolved, match)
			// full match may be resolved to non-string value in subclasses
			if len(matched) == 1 && len(regexMatch.Expr()) == len(resolved) {
				candidate := this.resolve(regexMatch)
				switch v := candidate.(type) {
				case string:
					sb.WriteString(v)
				default:
					return v
				}
			} else {
				sb.WriteString(resolved[lastMatch[1]:match[0]])
				sb.WriteString(fmt.Sprintf("%v", this.resolve(regexMatch)))
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

func (this *PatternProcessor) OverrideResolve(f func(*Match,
	func(*Match) any) any) {

	super := this.resolve
	this.resolve = func(rm *Match) any { return f(rm, super) }
}
