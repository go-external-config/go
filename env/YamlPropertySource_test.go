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
		source.ResolvePlaceholders()

		require.Equal(t, "value1", source.Property("a.key1"))
		require.Equal(t, 2.5, source.Property("a.key2"))
		require.Equal(t, 5, source.Property("ab.key1"))
		require.Equal(t, "h", source.Property("ab.key2"))
		require.Equal(t, "5", source.Property("c.array[0].value"))
		require.Equal(t, "7.5", source.Property("c.array[1].sub-array[0].sub1"))
	})
}
