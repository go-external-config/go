package util_test

import (
	"testing"

	"github.com/madamovych/go/util"
	"github.com/stretchr/testify/require"
)

func Test_Optional(t *testing.T) {
	t.Run("should produce expected results'", func(t *testing.T) {
		require.Equal(t, true, util.OptionalOfNilable(0).Present())
		require.Equal(t, true, util.OptionalOfNilable("").Present())

		var notInitializedAny any
		require.Equal(t, false, util.OptionalOfNilable(notInitializedAny).Present())

		var notInitializedPtr *any
		require.Equal(t, false, util.OptionalOfNilable(notInitializedPtr).Present())

		var notInitializedInt int
		require.Equal(t, true, util.OptionalOfNilable(notInitializedInt).Present())

		var notInitializedMap map[string]any
		require.Equal(t, false, util.OptionalOfNilable(notInitializedMap).Present())

		var notInitializedString string
		require.Equal(t, true, util.OptionalOfNilable(notInitializedString).Present())
	})
}
