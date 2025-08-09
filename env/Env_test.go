package env_test

import (
	"testing"

	"github.com/go-external-config/go/env"
	"github.com/stretchr/testify/require"
)

func Test_Env_Value(t *testing.T) {
	t.Run("should decode property", func(t *testing.T) {
		env.SetActiveProfiles("")
		env.Instance().AddPropertySource(env.MapPropertySourceOfMap("first loaded", map[string]string{
			"key": "value1"}))
		env.Instance().AddPropertySource(env.MapPropertySourceOfMap("second loaded", map[string]string{
			"key": "value2"}))
		env.Instance().AddPropertySource(env.MapPropertySourceOfMap("third loaded", map[string]string{
			"key":     "value3",
			"int":     "123",
			"intExpr": "#{123}",
			"servers": "host1,host2,host3",
			"slice":   "#{split('prod,live', ',')}"}))

		// last wins
		require.Equal(t, "value3", env.Value[string]("${key}"))

		// conversions
		require.Equal(t, "123", env.Value[string]("${int}"))
		require.Equal(t, 123, env.Value[int]("${int}"))
		require.Equal(t, float32(123.0), env.Value[float32]("${int}"))
		require.Equal(t, 123.0, env.Value[float64]("${int}"))
		require.Equal(t, "123", env.Value[string]("${intExpr}"))
		require.Equal(t, 123, env.Value[int]("${intExpr}"))
		require.Equal(t, float32(123.0), env.Value[float32]("${intExpr}"))
		require.Equal(t, "host1,host2,host3", env.Value[string]("${servers}"))
		require.Equal(t, []string{"host1", "host2", "host3"}, env.Value[[]string]("#{split('${servers}', ',')}"))
		require.Equal(t, []string{"prod", "live"}, env.Value[[]string]("${slice}"))
		require.Equal(t, "[prod live]", env.Value[string]("${slice}"))
		type Port int
		require.Equal(t, Port(123), env.Value[Port]("${intExpr}"))
	})
}

func Test_Env_ConfigurationProperties(t *testing.T) {
	t.Run("should decode property", func(t *testing.T) {
		env.SetActiveProfiles("")
		env.Instance().AddPropertySource(env.MapPropertySourceOfMap("properties", map[string]string{
			"key":      "value",
			"db.alias": "alias",
			"db.host":  "localhost",
			"db.port1": "111",
			"db.port3": "333"}))

		type Port int
		var db struct {
			Host  string
			host  string
			port1 int
			port2 int
			port3 Port
		}

		env.ConfigurationProperties("db", &db)

		require.Equal(t, "localhost", db.Host)
		require.Equal(t, "localhost", db.host)
		require.Equal(t, 111, db.port1)
		require.Equal(t, 0, db.port2)
		require.Equal(t, Port(333), db.port3)
	})
}

func Test_Env_MatchesProfiles(t *testing.T) {
	t.Run("should match profiles properly", func(t *testing.T) {
		env.SetActiveProfiles("test,hsqldb")

		require.True(t, env.MatchesProfiles())
		require.True(t, env.MatchesProfiles("test"))
		require.True(t, env.MatchesProfiles("hsqldb"))
		require.True(t, env.MatchesProfiles("!prod"))
		require.False(t, env.MatchesProfiles("prod"))
		require.False(t, env.MatchesProfiles("qa"))

		require.True(t, env.MatchesProfiles("test & hsqldb"))
		require.True(t, env.MatchesProfiles("(test & hsqldb) | prod"))
		require.False(t, env.MatchesProfiles("prod & hsqldb"))
		require.True(t, env.MatchesProfiles("prod | hsqldb"))
		require.False(t, env.MatchesProfiles("prod | qa"))

		require.True(t, env.MatchesProfiles("prod", "hsqldb"))
	})
}
