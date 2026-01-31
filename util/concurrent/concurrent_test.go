package concurrent_test

import (
	"sync"
	"testing"

	"github.com/go-external-config/go/util/concurrent"
	"github.com/stretchr/testify/require"
)

func Test_Synchronized(t *testing.T) {
	t.Run("should produce sane results", func(t *testing.T) {
		var x int
		var wg sync.WaitGroup
		var m sync.Mutex

		const n = 1000000
		wg.Add(n)

		for i := 0; i < n; i++ {
			go func() {
				defer wg.Done()
				concurrent.Synchronized(&m, func() {
					x++
				})
			}()
		}

		wg.Wait()
		require.Equal(t, n, x)
	})
}
