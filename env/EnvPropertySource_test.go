package env

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnvPropertySource(test *testing.T) {
	source := &EnvPropertySource{
		MapPropertySource: MapPropertySourceOf("Environment variables"),
	}

	source.SetProperty("SERVER_PORT", "9000")
	source.SetProperty("server.port", "8000")
	source.SetProperty("MAIN_LOGSTARTUPINFO", "true")
	source.SetProperty("MY_SERVICE_0__OTHER", "value")

	require.Equal(test, "8000", source.Property("server.port"))
	require.Equal(test, "9000", source.Property("SERVER_PORT"))
	require.Equal(test, "true", source.Property("main.log-startup-info"))
	require.Equal(test, "value", source.Property("my.service[0].other"))

	require.True(test, source.HasProperty("main.log-startup-info"))
	require.False(test, source.HasProperty("missing.property"))
	require.Panics(test, func() {
		source.Property("missing.property")
	})
}
