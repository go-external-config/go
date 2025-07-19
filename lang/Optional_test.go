package lang_test

import (
	"testing"

	"github.com/madamovych/go/lang"
	"github.com/stretchr/testify/require"
)

func Test_Optional(t *testing.T) {
	t.Run("should produce expected results'", func(t *testing.T) {
		require.Equal(t, true, lang.OptionalOfNilable(0).Present())

		var notInitializedAny any
		require.Equal(t, false, lang.OptionalOfNilable(notInitializedAny).Present())

		var notInitializedPtr *any
		require.Equal(t, false, lang.OptionalOfNilable(notInitializedPtr).Present())

		var notInitializedInt int
		require.Equal(t, true, lang.OptionalOfNilable(notInitializedInt).Present())

		var notInitializedMap map[string]any
		require.Equal(t, false, lang.OptionalOfNilable(notInitializedMap).Present())

		var notInitializedString string
		require.Equal(t, true, lang.OptionalOfNilable(notInitializedString).Present())
	})
}
