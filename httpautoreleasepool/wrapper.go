package httpautoreleasepool

import (
	"context"
	"net/http"

	"github.com/demdxx/autoreleasepool"
)

// Wrap wraps http.HandlerFunc with autorelease pool
func Wrap(hf http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		autoreleasepool.WithContext(r.Context(), func(ctx context.Context) error {
			hf(w, r.WithContext(ctx))
			return nil
		})
	}
}
