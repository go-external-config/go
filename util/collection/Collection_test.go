package collection_test

import (
	"testing"

	"github.com/madamovych/go/util/collection"
	"github.com/stretchr/testify/require"
)

func Test_Collection_SubtractSlice(t *testing.T) {
	t.Run("should produce sane results'", func(t *testing.T) {
		a := []int{1, 2, 3, 4, 5}
		b := []int{4, 3, 2}
		require.Equal(t, []int{1, 5}, collection.SubtractSlice(a, b))
	})
}

func Test_Collection_ReverseSlice(t *testing.T) {
	t.Run("should produce sane results'", func(t *testing.T) {
		a := []int{1, 2, 3, 4, 5}
		require.Equal(t, []int{5, 4, 3, 2, 1}, collection.ReverseSlice(a))
	})
}
