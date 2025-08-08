package regex

import (
	"regexp"

	"github.com/go-external-config/v1/util"
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

func (m *Match) Expr() string {
	return m.Group(0).Value()
}

func (m *Match) GroupCount() int {
	return len(m.subexpByIdx)
}

func (m *Match) Group(idx int) *util.Optional[string] {
	return util.OptionalOfEntry(m.subexpByIdx, idx)
}

func (m *Match) NamedGroup(name string) *util.Optional[string] {
	return util.OptionalOfEntry(m.subexpByName, name)
}
