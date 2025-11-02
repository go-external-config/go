package env_test

import (
	"math"
	"testing"

	"github.com/go-external-config/go/env"
	"github.com/stretchr/testify/require"
)

func Test_RandomValues(t *testing.T) {
	require.Equal(t, 36, len(env.Value[string]("${random.uuid}")))

	require.Equal(t, 10, len(env.Value[string]("${random.string(10)}")))

	require.Equal(t, 8, len(env.Value[string]("${random.value(4)}")))
	require.Equal(t, 16, len(env.Value[string]("${random.value}")))
	require.Equal(t, 32, len(env.Value[string]("${random.value(16)}")))

	for i := 0; i < 1000; i++ {
		require.LessOrEqual(t, math.MinInt32, env.Value[int]("${random.int}"))
		require.Greater(t, math.MaxInt32, env.Value[int]("${random.int}"))
	}
	for i := 0; i < 1000; i++ {
		require.LessOrEqual(t, 0, env.Value[int]("${random.int(10)}"))
		require.Greater(t, 10, env.Value[int]("${random.int(10)}"))
	}
	for i := 0; i < 1000; i++ {
		require.LessOrEqual(t, -10, env.Value[int]("${random.int(-10,10)}"))
		require.Greater(t, 10, env.Value[int]("${random.int(-10,10)}"))
	}

	for i := 0; i < 1000; i++ {
		require.LessOrEqual(t, math.MinInt64, env.Value[int]("${random.int64}"))
		require.Greater(t, math.MaxInt64, env.Value[int]("${random.int64}"))
	}
	for i := 0; i < 1000; i++ {
		require.LessOrEqual(t, 0, env.Value[int]("${random.int64(10)}"))
		require.Greater(t, 10, env.Value[int]("${random.int64(10)}"))
	}
	for i := 0; i < 1000; i++ {
		require.LessOrEqual(t, -10, env.Value[int]("${random.int64(-10,10)}"))
		require.Greater(t, 10, env.Value[int]("${random.int64(-10,10)}"))
	}
}
