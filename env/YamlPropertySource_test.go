package env_test

import (
	"testing"

	"github.com/madamovych/go/env"
	"github.com/stretchr/testify/require"
)

func Test_YamlPropertySource_Resolve(t *testing.T) {
	t.Run("should resolve variables", func(t *testing.T) {
		source := env.NewYamlPropertySource("yamlPropertySource", `
a:
  key1: "value1"
  key2: 2.5
ab:
  key1: 5
  key2: "h"
b:
  emptyKey:
c:
  array:
    - name: "element1"
      value: ${prop5}
    - name: "element2"
      sub-array:
        - sub1: "#{${a.key2} + ${ab.key1}}"
`)
		source.SetProperty("prop5", "#{3+2}")
		processor := env.ExprProcessorOf(true)
		processor.SetPropertySource(source)

		require.Equal(t, "value1", processor.Process("${a.key1}"))
		require.Equal(t, "2.5", processor.Process("${a.key2}"))
		require.Equal(t, "5", processor.Process("${ab.key1}"))
		require.Equal(t, "h", processor.Process("${ab.key2}"))
		require.Equal(t, "5", processor.Process("${c.array[0].value}"))
		require.Equal(t, "7.5", processor.Process("${c.array[1].sub-array[0].sub1}"))
	})
}
