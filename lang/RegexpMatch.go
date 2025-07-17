package lang

import "regexp"

type RegexpMatch struct {
	expr         string
	subexpByIdx  []string
	subexpByName map[string]string
}

func RegexpMatchOf(regexp *regexp.Regexp, str string, matched []int) *RegexpMatch {
	subexpCount := len(regexp.SubexpNames())
	result := RegexpMatch{
		expr:         str[matched[0]:matched[1]],
		subexpByIdx:  make([]string, subexpCount),
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

func (m *RegexpMatch) Expr() string {
	return m.Group(0)
}

func (m *RegexpMatch) GroupCount() int {
	return len(m.subexpByIdx)
}

func (m *RegexpMatch) Group(idx int) string {
	return m.subexpByIdx[idx]
}

func (m *RegexpMatch) NamedGroup(name string) string {
	return m.subexpByName[name]
}
