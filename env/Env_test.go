package env_test

import (
	"testing"

	"github.com/madamovych/go/env"
	"github.com/stretchr/testify/require"
)

func Test_Env_Value(t *testing.T) {
	t.Run("should decode property", func(t *testing.T) {
		env.SetActiveProfiles("test")
		env.GetEnvironment().AddPropertySource(env.MapPropertySourceOfMap("first loaded", map[string]string{
			"key": "value1"}))
		env.GetEnvironment().AddPropertySource(env.MapPropertySourceOfMap("second loaded", map[string]string{
			"key": "value2"}))
		env.GetEnvironment().AddPropertySource(env.MapPropertySourceOfMap("third loaded", map[string]string{
			"key":     "value3",
			"int":     "123",
			"int8":    "123",
			"int16":   "123",
			"int32":   "123",
			"int64":   "123",
			"uint":    "123",
			"uint8":   "123",
			"uint16":  "123",
			"uint32":  "123",
			"uint64":  "123",
			"float32": "123",
			"float64": "123",
			"bool":    "true",
			"string":  "123"}))

		require.Equal(t, "value3", env.Value[string]("${key}"))

		var v_int int = env.Value[int]("${int}")
		require.Equal(t, 123, v_int)
		var v_int8 int8 = env.Value[int8]("${int8}")
		require.Equal(t, int8(123), v_int8)
		var v_int16 int16 = env.Value[int16]("${int16}")
		require.Equal(t, int16(123), v_int16)
		var v_int32 int32 = env.Value[int32]("${int32}")
		require.Equal(t, int32(123), v_int32)
		var v_int64 int64 = env.Value[int64]("${int64}")
		require.Equal(t, int64(123), v_int64)
		var v_uint uint = env.Value[uint]("${uint}")
		require.Equal(t, uint(123), v_uint)
		var v_uint8 uint8 = env.Value[uint8]("${uint8}")
		require.Equal(t, uint8(123), v_uint8)
		var v_uint16 uint16 = env.Value[uint16]("${uint16}")
		require.Equal(t, uint16(123), v_uint16)
		var v_uint32 uint32 = env.Value[uint32]("${uint32}")
		require.Equal(t, uint32(123), v_uint32)
		var v_uint64 uint64 = env.Value[uint64]("${uint64}")
		require.Equal(t, uint64(123), v_uint64)
		var v_float32 float32 = env.Value[float32]("${float32}")
		require.Equal(t, float32(123), v_float32)
		var v_float64 float64 = env.Value[float64]("${float64}")
		require.Equal(t, float64(123), v_float64)
		var v_bool bool = env.Value[bool]("${bool}")
		require.Equal(t, true, v_bool)
		var v_string string = env.Value[string]("${string}")
		require.Equal(t, "123", v_string)
	})
}
