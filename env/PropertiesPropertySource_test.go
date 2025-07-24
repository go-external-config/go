package env_test

import (
	"testing"

	"github.com/madamovych/go/env"
	"github.com/stretchr/testify/require"
)

func Test_PropertiesPropertySource_Resolve(t *testing.T) {
	t.Run("should resolve variables", func(t *testing.T) {
		source := env.NewPropertiesPropertySource("mapPropertySource", `
		prop1=val1
		prop2=${prop1}
		prop3 = #{'${prop1}'}
		prop4=${prop#{${prop5}-2}}
`)
		source.SetProperty("prop5", "#{2+2}")
		processor := env.ExprProcessorOf(true)
		processor.SetPropertySource(source)

		require.Equal(t, "val1", processor.Process("${prop1}"))
		require.Equal(t, "val1", processor.Process("${prop2}"))
		require.Equal(t, "val1", processor.Process("${prop3}"))
		require.Equal(t, "val1", processor.Process("${prop4}"))
		require.Equal(t, "4", processor.Process("${prop5}"))
	})
}
