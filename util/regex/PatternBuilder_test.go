package regex_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/go-external-config/go/util/regex"
	"github.com/stretchr/testify/require"
)

func TestPatternBuilder_dummy(test *testing.T) {
	str := "Input param ABC from block4 as	abc \n next line"
	pattern := regex.NewPatternBuilder().Next(`param {paramName} from (B|b)lock{blockNum:\d}`).Opt(` as {varName:\w+}`).Build()
	compiled := regexp.MustCompile(pattern)
	for _, m := range compiled.FindAllStringSubmatchIndex(str, -1) {
		match := regex.MatchOf(compiled, str, m)
		require.Equal(test, "ABC", match.NamedGroup("paramName").Value())
		require.Equal(test, "4", match.NamedGroup("blockNum").Value())
		require.Equal(test, "abc", match.NamedGroup("varName").Value())
		return
	}
	require.Fail(test, "Doesn't match")
}

func TestPatternBuilder_modifier(test *testing.T) {
	str := "Input param ABC from block4 as	abc \n next line"
	pattern := regex.NewPatternBuilder(`(?m)`).Next(`param {paramName} from (B|b)lock{blockNum:\d{1}}`).Opt(` as {varName}`).End().Build()
	fmt.Println(pattern)
	compiled := regexp.MustCompile(pattern)
	for _, m := range compiled.FindAllStringSubmatchIndex(str, -1) {
		match := regex.MatchOf(compiled, str, m)
		require.Equal(test, "ABC", match.NamedGroup("paramName").Value())
		require.Equal(test, "4", match.NamedGroup("blockNum").Value())
		require.Equal(test, "abc", match.NamedGroup("varName").Value())
		return
	}
	require.Fail(test, "Doesn't match")
}

func TestPatternBuilder_multiline(test *testing.T) {
	str := `
	Table ABC {trim=false} abc
	|name |rank    |
	|Larry|Stooge 3| 
	|Moe  |Stooge 1|
	|Curly|Stooge 2|
	
	something else
	`
	pattern := regex.NewPatternBuilder().
		Start(`Table {name}`).Opt(` \{{params}\}`).Next(" {alias}").End().
		Start(`{header}`).End().
		Start(`{body}`).
		Start(``).End().Build()

	compiled := regexp.MustCompile(pattern)
	for _, m := range compiled.FindAllStringSubmatchIndex(str, -1) {
		match := regex.MatchOf(compiled, str, m)
		require.Equal(test, "ABC", match.NamedGroup("name").Value())
		require.Equal(test, "trim=false", match.NamedGroup("params").Value())
		require.Equal(test, "abc", match.NamedGroup("alias").Value())
		require.Equal(test, "	|name |rank    |", match.NamedGroup("header").Value())
		body := `	|Larry|Stooge 3| 
	|Moe  |Stooge 1|
	|Curly|Stooge 2|
`
		require.Equal(test, body, match.NamedGroup("body").Value())
		return
	}
	require.Fail(test, "Doesn't match")
}

func TestPatternBuilder_misc(test *testing.T) {
	require.Equal(test, `(?ms)(?P<location>.+)\[(?P<fantomExt>\.[\w]+)\]`, regex.NewPatternBuilder().Next(`{location:.+}\[{fantomExt:\.[\w]+}\]`).Build())
	require.Equal(test, "(?ms)(?P<word>\\w+)|(?P<sign>\\W)", regex.NewPatternBuilder().Next("{word:\\w+}|{sign:\\W}").Build())
	require.Equal(test, `(?ms)(?P<key>[^=\s]+)=(?P<value>.*)`, regex.NewPatternBuilder().Next(`{key:[^=\s]+}={value:.*}`).Build())
	require.Equal(test, `(?ms)--?(?P<key>[^=\s]+)\s*=?(?P<value>.*)`, regex.NewPatternBuilder().Next(`--?{key:[^=\s]+}\s*=?{value:.*}`).Build())
	require.Equal(test, `(?ms)^random\.((?P<uuid>uuid)|(?P<string>string)\((?P<size>\d+)\)|(?P<value>value)(\((?P<bytes>\d+)\))?|(?P<int>int)(\((?P<max>\d+)\))?|(?P<int>int)(\((?P<min>-?\d+),(?P<max>\d+)\))?|(?P<int64>int64)(\((?P<max>\d+)\))?|(?P<int64>int64)(\((?P<min>-?\d+),(?P<max>\d+)\))?)$`,
		regex.NewPatternBuilder().Next(`^random\.({uuid:uuid}|{string:string}\({size:\d+}\)|{value:value}(\({bytes:\d+}\))?|{int:int}(\({max:\d+}\))?|{int:int}(\({min:-?\d+},{max:\d+}\))?|{int64:int64}(\({max:\d+}\))?|{int64:int64}(\({min:-?\d+},{max:\d+}\))?)$`).Build())
}
