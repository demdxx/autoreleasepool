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

func BenchmarkAp(b *testing.B) {
	const sliceSize = 1000
	b.ReportAllocs()
	b.ResetTimer()

	b.Run("Do", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := Do(testDo(b))
			assert.NoError(b, err)
		}
	})

	b.Run("WithoutAuthoreleasepool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v := new(int)
			*v = 42
			a := make([]int, sliceSize, sliceSize)
			for i := 0; i < sliceSize; i++ {
				a[i] = *v
			}
			b := make([]int, sliceSize, sliceSize)
			for i := 0; i < sliceSize; i++ {
				b[i] = a[i]
			}
			b[0] = 0
		}
	})

	b.Run("WithAuthoreleasepool", func(b *testing.B) {
		_ = Do(func(ctx context.Context) error {
			for i := 0; i < b.N; i++ {
				v := New[int](ctx)
				*v = 42
				a := MakeSlice[int](ctx, sliceSize, sliceSize)
				for i := 0; i < sliceSize; i++ {
					a[i] = *v
				}
				b := MakeSlice[int](ctx, sliceSize, sliceSize)
				for i := 0; i < sliceSize; i++ {
					b[i] = a[i]
				}
				b[0] = 0
			}
			return nil
		})
	})

	b.Run("WithAuthoreleasepool2", func(b *testing.B) {
		_ = Do(func(ctx context.Context) error {
			for i := 0; i < b.N; i++ {
				v := New[int](ctx)
				*v = 42
				a := MakeSlice[int](ctx, sliceSize, sliceSize)
				for i := 0; i < sliceSize; i++ {
					a[i] = *v
				}
				b := MakeSlice[int](ctx, sliceSize, sliceSize)
				for i := 0; i < sliceSize; i++ {
					b[i] = a[i]
				}
				b[0] = 0
			}
			return nil
		})
	})
}
