package optional_test

import (
	"testing"

	"github.com/go-external-config/go/util/optional"
	"github.com/stretchr/testify/require"
)

func Test_Optional(t *testing.T) {
	t.Run("should produce expected results'", func(t *testing.T) {
		require.Equal(t, true, optional.OfNilable(0).Present())
		require.Equal(t, true, optional.OfNilable("").Present())

		var notInitializedAny any
		require.Equal(t, false, optional.OfNilable(notInitializedAny).Present())

		var notInitializedPtr *any
		require.Equal(t, false, optional.OfNilable(notInitializedPtr).Present())

		var notInitializedInt int
		require.Equal(t, true, optional.OfNilable(notInitializedInt).Present())

		var notInitializedMap map[string]any
		require.Equal(t, false, optional.OfNilable(notInitializedMap).Present())

		var notInitializedString string
		require.Equal(t, true, optional.OfNilable(notInitializedString).Present())
	})
}
