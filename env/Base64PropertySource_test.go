package env_test

import (
	"testing"

	"github.com/go-external-config/go/env"
	"github.com/stretchr/testify/require"
)

func Test_Base64PropertySource(t *testing.T) {
	t.Run("should decode property", func(t *testing.T) {
		env.SetActiveProfiles("")
		env.Instance().AddPropertySource(env.MapPropertySourceOfMap("delegate", map[string]string{
			"base64Encoded": "base64:IEhlbGxvIFdvcmxkISA="}))
		env.Instance().AddPropertySource(env.NewBase64PropertySource(env.Instance()))
		require.Equal(t, " Hello World! ", env.Value[string]("${base64Encoded}"))
	})
}
