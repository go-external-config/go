package env

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnvironmentPropertySourcePrecedence(test *testing.T) {
	config := MapPropertySourceOf("Configuration")
	config.SetProperty("server.port", "7000")

	environ := &EnvPropertySource{
		MapPropertySource: MapPropertySourceOf("Environment variables"),
	}
	environ.SetProperty("SERVER_PORT", "8000")

	params := &ParamsPropertySource{
		MapPropertySource: MapPropertySourceOf("Application parameters"),
	}
	params.SetProperty("server.port", "9000")

	environment := Environment{
		sources:             []PropertySource{config, environ, params},
		strictExprProcessor: NewExprProcessor(true),
	}

	require.Equal(test, "9000", environment.property("server.port"))

	environment.sources = []PropertySource{config, environ}
	require.Equal(test, "8000", environment.property("server.port"))

	environment.sources = []PropertySource{config}
	require.Equal(test, "7000", environment.property("server.port"))
}
