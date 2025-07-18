package env

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_PropertySource_Resolve(t *testing.T) {
	t.Run("should resolve variables", func(t *testing.T) {
		source := PropertySourceOfMap("mapPropertySource", map[string]any{
			"prop1": "val1",
			"prop2": "val2",
			"prop3": "#{val3${prop4}}",
			"prop4": "val#{2+2}"})
		source.SetProperty("prop5", "val5")

		source.ResolvePlaceholders()

		require.Equal(t, "val1", source.Property("prop1"))
		require.Equal(t, "val2", source.Property("prop2"))
		require.Equal(t, "val34", source.Property("prop3"))
		require.Equal(t, "val4", source.Property("prop4"))
		require.Equal(t, "val5", source.Property("prop5"))
	})
}
