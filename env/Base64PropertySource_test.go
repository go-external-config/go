package env_test

import (
	"testing"

	"github.com/madamovych/go/env"
	"github.com/stretchr/testify/require"
)

func Test_Base64PropertySourceTest(t *testing.T) {
	t.Run("should decode property", func(t *testing.T) {
		env.ActiveteProfiles("test")
		env.GetEnvironment().AddPropertySource(env.MapPropertySourceOfMap("delegate", map[string]string{
			"base64Encoded": "Base64:IEhlbGxvIFdvcmxkISA="}))
		env.GetEnvironment().AddPropertySource(env.NewBase64PropertySource())
		require.Equal(t, " Hello World! ", env.Value("${base64Encoded}"))
	})
}
