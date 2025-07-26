package env_test

import (
	"testing"

	"github.com/go-external-config/go/env"
	"github.com/stretchr/testify/require"
)

func Test_Env_PropertyVsValue(t *testing.T) {
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
		require.Equal(t, "value3", env.Property("key"))
		require.Equal(t, "value3", env.Value("${key}"))

		// property values are strings (despite yaml allow other types) unless property value is an expression
		require.Equal(t, "123", env.Property("int"))
		require.Equal(t, "123", env.Value("${int}"))

		// string values can be converted to compatible types
		require.Equal(t, 123, env.PropertyAs[int]("int"))
		require.Equal(t, 123, env.ValueAs[int]("${int}"))
		require.Equal(t, float32(123.0), env.PropertyAs[float32]("int"))
		require.Equal(t, 123.0, env.PropertyAs[float64]("int"))

		// expression evaluation can produce any types
		require.Equal(t, 123, env.Property("intExpr"))
		require.Equal(t, 123, env.Value("${intExpr}"))

		// any types (resulted from expression evaluation) can be converted to compatible types
		require.Equal(t, 123, env.PropertyAs[int]("intExpr"))
		require.Equal(t, float32(123.0), env.PropertyAs[float32]("intExpr"))
		require.Equal(t, "123", env.PropertyAs[string]("intExpr"))

		require.Equal(t, "host1,host2,host3", env.Property("servers"))
		require.Equal(t, "host1,host2,host3", env.Value("${servers}"))
		require.Equal(t, []string{"host1", "host2", "host3"}, env.Value("#{split('${servers}', ',')}"))

		require.Equal(t, []string{"prod", "live"}, env.Property("slice"))
		require.Equal(t, []string{"prod", "live"}, env.Value("${slice}"))
		require.Equal(t, "[prod live]", env.PropertyAs[string]("slice"))

		type Port int
		require.Equal(t, Port(123), env.PropertyAs[Port]("intExpr"))
	})
}

func Test_Env_ConfigurationProperties(t *testing.T) {
	t.Run("should decode property", func(t *testing.T) {
		env.SetActiveProfiles("")
		env.Instance().AddPropertySource(env.MapPropertySourceOfMap("properties", map[string]string{
			"key":      "value",
			"db.alias": "alias",
			"db.Host":  "localhost",
			"db.port1": "111",
			"db.port3": "333"}))

		type Port int
		var db struct {
			Host  string
			port1 int
			port2 int
			port3 Port
		}

		env.ConfigurationProperties("db", &db)

		require.Equal(t, "localhost", db.Host)
		require.Equal(t, 111, db.port1)
		require.Equal(t, 0, db.port2)
		require.Equal(t, Port(333), db.port3)
	})
}
