package text_test

import (
	"testing"

	"github.com/madamovych/go/lang"
	"github.com/madamovych/go/text"
	"github.com/stretchr/testify/require"
)

func Test_PatternProcessor_Process(t *testing.T) {
	t.Run("should substitute 'World'", func(t *testing.T) {
		processor := text.NewPatternProcessor("\\w+")
		processor.SetResolve(func(match *lang.RegexpMatch) string {
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
