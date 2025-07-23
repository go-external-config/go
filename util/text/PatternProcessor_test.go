package text_test

import (
	"testing"

	"github.com/madamovych/go/util/regex"
	"github.com/madamovych/go/util/text"
	"github.com/stretchr/testify/require"
)

func Test_PatternProcessor_Process(t *testing.T) {
	t.Run("should substitute 'World'", func(t *testing.T) {
		processor := text.PatternProcessorOf("\\w+")
		processor.OverrideResolve(func(match *regex.Match,
			super func(*regex.Match) string) string {

			switch match.Expr() {
			case "World":
				return "Mike"
			default:
				return match.Expr()
			}
		})

		require.Equal(t, " Hello Mike! ", processor.Process(" Hello World! "))
	})
}
