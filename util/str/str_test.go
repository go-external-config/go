package str_test

import (
	"testing"

	"github.com/go-external-config/go/util/str"
	"github.com/stretchr/testify/require"
)

func Test_Parse(t *testing.T) {
	t.Run("parse string as different values'", func(t *testing.T) {
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

func Test_ReplaceChars(t *testing.T) {
	t.Run("replace characters in a string'", func(t *testing.T) {
		rules := map[rune]rune{
			'.': '_',
			'[': '_',
			']': '_',
			'-': 0, // delete
		}
		require.Equal(t, "FOO_BARBAZ", str.ReplaceChars("foo.bar-baz", rules))
		require.Equal(t, "A_B_CD", str.ReplaceChars("a[b]c-d", rules))
	})
}

func TestJoin(t *testing.T) {

	// No values
	require.Equal(t, "", str.Join(","))

	// Single values
	require.Equal(t, "a", str.Join(",", "a"))
	require.Equal(t, "", str.Join(",", ""))

	// Multiple values, no empties
	require.Equal(t, "a,b,c", str.Join(",", "a", "b", "c"))
	require.Equal(t, "a,b", str.Join(",", "", "a", "b"))
	require.Equal(t, "a,b", str.Join(",", "a", "", "b"))
	require.Equal(t, "a,b,", str.Join(",", "a", "b", ""))
	require.Equal(t, "a,b", str.Join(",", "a", "", "", "b"))
	require.Equal(t, "ab", str.Join("", "a", "b"))
	require.Equal(t, "ab", str.Join("", "", "a", "b", ""))
	require.Equal(t, "1-2-3", str.Join("-", 1, 2, 3))
	require.Equal(t, "x 5 true", str.Join(" ", "x", 5, true))

	// nil values
	require.Equal(t, "<nil>", str.Join(",", nil))
	require.Equal(t, "a,<nil>,b", str.Join(",", "a", nil, "b"))
}
