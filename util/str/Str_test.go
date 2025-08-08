package str_test

import (
	"testing"

	"github.com/go-external-config/v1/util/str"
	"github.com/stretchr/testify/require"
)

func Test_Parse(t *testing.T) {
	t.Run("should produce sane results'", func(t *testing.T) {
		value := "123"

		require.Equal(t, int(123), str.Parse[int](value))
		require.Equal(t, int8(123), str.Parse[int8](value))
		require.Equal(t, int16(123), str.Parse[int16](value))
		require.Equal(t, int32(123), str.Parse[int32](value))
		require.Equal(t, int64(123), str.Parse[int64](value))
		require.Equal(t, uint(123), str.Parse[uint](value))
		require.Equal(t, uint8(123), str.Parse[uint8](value))
		require.Equal(t, uint16(123), str.Parse[uint16](value))
		require.Equal(t, uint32(123), str.Parse[uint32](value))
		require.Equal(t, uint64(123), str.Parse[uint64](value))
		require.Equal(t, float32(123), str.Parse[float32](value))
		require.Equal(t, float64(123), str.Parse[float64](value))
		require.Equal(t, "123", str.Parse[string](value))
		type Port int8
		require.Equal(t, Port(123), str.Parse[Port](value))
		require.Equal(t, true, str.Parse[bool]("true"))
	})
}
