//go:build goexperiment.arenas

package autoreleasepool

import (
	"arena"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testDo(t assert.TestingT) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		v := New[int](ctx)
		*v = 42
		assert.Equal(t, 42, *v)

		a := MakeSlice[int](ctx, 10, 10)
		for i := 0; i < 10; i++ {
			a[i] = i
		}
		assert.ElementsMatch(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, a)

		return nil
	}
}

func TestDo(t *testing.T) {
	t.Run("Do", func(t *testing.T) {
		err := Do(testDo(t))
		assert.NoError(t, err)
	})

	t.Run("WithContext", func(t *testing.T) {
		err := WithContext(context.Background(), testDo(t))
		assert.NoError(t, err)
	})

	t.Run("WithArena", func(t *testing.T) {
		mem := arena.NewArena()
		defer mem.Free()
		err := WithArena(mem, testDo(t))
		assert.NoError(t, err)
	})

	t.Run("WithContextAndArena", func(t *testing.T) {
		mem := arena.NewArena()
		defer mem.Free()
		err := WithContextAndArena(context.Background(), mem, testDo(t))
		assert.NoError(t, err)
	})
}

func BenchmarkDo(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := Do(testDo(b))
			assert.NoError(b, err)
		}
	})
}
