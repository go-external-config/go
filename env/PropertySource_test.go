package env

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_PropertySource_Resolve(t *testing.T) {
	t.Run("should resolve variables", func(t *testing.T) {
		source := PropertySourceOfMap("mapPropertySource", map[string]any{
			"prop1": "val1",
			"prop2": "${prop1}",
			"prop3": "#{prop1}",
			"prop4": "${prop#{${prop5}-2}}"})
		source.SetProperty("prop5", "#{2+2}")
		source.ResolvePlaceholders()

		require.Equal(t, "val1", source.Property("prop1"))
		require.Equal(t, "val1", source.Property("prop2"))
		require.Equal(t, "val1", source.Property("prop3"))
		require.Equal(t, "val1", source.Property("prop4"))
		require.Equal(t, "4", source.Property("prop5"))
	})
}
