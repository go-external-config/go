package env_test

import (
	"testing"

	"github.com/go-external-config/go/env"
	"github.com/stretchr/testify/require"
)

func Test_MapPropertySource_Resolve(t *testing.T) {
	t.Run("should resolve variables", func(t *testing.T) {
		source := env.MapPropertySourceOfMap("mapPropertySource", map[string]string{
			"prop1": "val1",
			"prop2": "${prop1}",
			"prop3": "#{'${prop1}'}",
			"prop4": "${prop#{${prop5}-2}}"})
		source.SetProperty("prop5", "#{2+2}")
		env.RegisterPropertySource(source)

		require.Equal(t, "val1", env.Property("prop1"))
		require.Equal(t, "val1", env.Property("prop2"))
		require.Equal(t, "val1", env.Property("prop3"))
		require.Equal(t, "val1", env.Property("prop4"))
		require.Equal(t, "4", env.Property("prop5"))
	})
}
