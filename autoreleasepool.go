//go:build goexperiment.arenas

package autoreleasepool

import (
	"arena"
	"context"
)

var ctxArenaKey = &struct{}{}

// WithContextAndArena runs a new autorelease pool and executes the given function
func WithContextAndArena(ctx context.Context, mem *arena.Arena, fn func(ctx context.Context) error) error {
	ctx = context.WithValue(ctx, ctxArenaKey, mem)
	return fn(ctx)
}

// WithContext runs a new autorelease pool and executes the given function
func WithContext(ctx context.Context, fn func(ctx context.Context) error) error {
	mem := arena.NewArena()
	defer mem.Free()

	return WithContextAndArena(ctx, mem, fn)
}

// WithArena runs a new autorelease pool and executes the given function
func WithArena(mem *arena.Arena, fn func(ctx context.Context) error) error {
	return WithContextAndArena(context.Background(), mem, fn)
}

// Do runs a new autorelease pool and executes the given function
func Do(fn func(ctx context.Context) error) error {
	return WithContext(context.Background(), fn)
}

// New allocates a new object of type T in the current autorelease pool
func New[T any](ctx context.Context) *T {
	are := ctx.Value(ctxArenaKey)
	if are == nil {
		return new(T)
	}
	mem := are.(*arena.Arena)
	return arena.New[T](mem)
}

// MakeSlice allocates a new slice of type T in the current autorelease pool
// In case of a nil context, it will call make() instead
func MakeSlice[T any](ctx context.Context, len, cap int) []T {
	are := ctx.Value(ctxArenaKey)
	if are == nil {
		return make([]T, len, cap)
	}
	mem := are.(*arena.Arena)
	return arena.MakeSlice[T](mem, len, cap)
}

// Clone clones the given object in the current autorelease pool
// In case of a nil context, it will return the original object
func Clone[T any](ctx context.Context, s T) T {
	are := ctx.Value(ctxArenaKey)
	if are == nil {
		return s
	}
	return arena.Clone[T](s)
}
