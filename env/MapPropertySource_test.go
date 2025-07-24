package env

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_MapPropertySource_Resolve(t *testing.T) {
	t.Run("should resolve variables", func(t *testing.T) {
		source := MapPropertySourceOfMap("mapPropertySource", map[string]string{
			"prop1": "val1",
			"prop2": "${prop1}",
			"prop3": "#{'${prop1}'}",
			"prop4": "${prop#{${prop5}-2}}"})
		source.SetProperty("prop5", "#{2+2}")
		processor := ExprProcessorOf(true)
		processor.SetPropertySource(source)

		require.Equal(t, "val1", processor.Process("${prop1}"))
		require.Equal(t, "val1", processor.Process("${prop2}"))
		require.Equal(t, "val1", processor.Process("${prop3}"))
		require.Equal(t, "val1", processor.Process("${prop4}"))
		require.Equal(t, "4", processor.Process("${prop5}"))
	})
}
