package regex_test

import (
	"regexp"
	"testing"

	"github.com/go-external-config/go/util/regex"
	"github.com/stretchr/testify/require"
)

func Test_Match(t *testing.T) {
	t.Run("should produce sane results'", func(t *testing.T) {
		pattern := regexp.MustCompile(`\$\$\{(?P<long>.*?)\}\$|\$\{(?P<short>.*?)\}`)
		str := `To be, or not to be: that is the ${question}:
			Whether 'tis nobler in the mind to suffer
			The slings and arrows of outrageous fortune,
			Or to take arms against a sea of troubles,
			And by opposing end them? To die: to sleep;`

		for _, m := range pattern.FindAllStringSubmatchIndex(str, -1) {
			match := regex.MatchOf(pattern, str, m)
			require.Equal(t, "${question}", match.Expr())
			require.Equal(t, "question", match.Group(2).Value())
			require.Equal(t, "question", match.NamedGroup("short").Value())
			require.Equal(t, false, match.Group(1).Present())
			require.Equal(t, true, match.NamedGroup("short").Present())
			require.Equal(t, false, match.NamedGroup("long").Present())
		}
	})
}
