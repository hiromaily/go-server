package middleware

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
)

const requestIDKey = key(42)

// GetRequestID is to get request ID
func GetRequestID(ctx context.Context) (int64, error) {
	id, ok := ctx.Value(requestIDKey).(int64)
	if !ok {
		return 0, fmt.Errorf("%s", "couldn't find request ID in context")
	}
	return id, nil
}

// SetRequestID is to set request ID
func SetRequestID() Handler {
	return func(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request) {
		ctx := r.Context()
		id := rand.Int63()
		ctx = context.WithValue(ctx, requestIDKey, id)

		return w, r.WithContext(ctx)
	}
}
