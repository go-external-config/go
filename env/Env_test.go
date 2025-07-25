package env_test

import (
	"fmt"
	"testing"

	"github.com/go-external-config/go/env"
	"github.com/stretchr/testify/require"
)

func Test_Env_Value(t *testing.T) {
	t.Run("should decode property", func(t *testing.T) {
		env.SetActiveProfiles("")
		env.EnvironmentInstance().AddPropertySource(env.MapPropertySourceOfMap("first loaded", map[string]string{
			"key": "value1"}))
		env.EnvironmentInstance().AddPropertySource(env.MapPropertySourceOfMap("second loaded", map[string]string{
			"key": "value2"}))
		env.EnvironmentInstance().AddPropertySource(env.MapPropertySourceOfMap("third loaded", map[string]string{
			"key":       "value3",
			"intString": "123",
			"int":       "#{123}",
			"arrString": "#{split('prod,live', ',')}"}))

		require.Equal(t, "value3", env.Value("${key}"))
		require.Equal(t, "123", env.Value("${intString}"))
		require.Equal(t, 123, env.Value("${int}"))
		require.Equal(t, []string{"prod", "live"}, env.Value("${arrString}"))
	})
}

func Test_Env_ConfigurationProperties(t *testing.T) {
	t.Run("should decode property", func(t *testing.T) {
		env.SetActiveProfiles("")
		env.EnvironmentInstance().AddPropertySource(env.MapPropertySourceOfMap("properties", map[string]string{
			"key":     "value",
			"db.Host": "localhost",
			"db.Port": "123"}))

		var db struct {
			Host string
			Port string
		}

		db = env.ConfigurationProperties("db", db)

		fmt.Printf("%s\n", db.Host)
		fmt.Printf("%s\n", db.Port)
	})
}
