package text

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProcess(t *testing.T) {
	t.Run("should substitute 'World'", func(t *testing.T) {
		processor := *NewtPatternProcessor("\\w+")
		processor.resolve = func(matcher string) string {
			switch matcher {
			case "World":
				return "Mike"
			default:
				return matcher
			}
		}

		require.Equal(t, " Hello Mike! ", processor.Process(" Hello World! "))
	})
}
