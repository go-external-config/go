package regex

import (
	"regexp"

	"github.com/go-external-config/go/util/optional"
)

type Match struct {
	expr         string
	subexpByIdx  map[int]string
	subexpByName map[string]string
}

func MatchOf(regexp *regexp.Regexp, str string, matched []int) *Match {
	subexpCount := len(regexp.SubexpNames())
	result := Match{
		expr:         str[matched[0]:matched[1]],
		subexpByIdx:  make(map[int]string, subexpCount),
		subexpByName: make(map[string]string, subexpCount)}

	for idx, name := range regexp.SubexpNames() {
		if matched[2*idx] >= 0 {
			value := str[matched[2*idx]:matched[2*idx+1]]
			result.subexpByIdx[idx] = value
			if name != "" {
				result.subexpByName[name] = value
			}
		}
	}
	return &result
}

func (this *Match) Expr() string {
	return this.Group(0).Value()
}

func (this *Match) GroupCount() int {
	return len(this.subexpByIdx)
}

func (this *Match) Group(idx int) *optional.Optional[string] {
	return optional.OfEntry(this.subexpByIdx, idx)
}

func (this *Match) NamedGroup(name string) *optional.Optional[string] {
	return optional.OfEntry(this.subexpByName, name)
}
