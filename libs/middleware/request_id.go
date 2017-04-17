package middleware

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
)

type key int

const requestIDKey = key(42)

func GetRequestID(ctx context.Context) (int64, error) {
	id, ok := ctx.Value(requestIDKey).(int64)
	if !ok {
		return 0, fmt.Errorf("%s", "couldn't find request ID in context")
	}
	return id, nil
}

func SetRequestID() Handler {
	return func(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request) {
		ctx := r.Context()
		id := rand.Int63()
		ctx = context.WithValue(ctx, requestIDKey, id)

		return w, r.WithContext(ctx)
	}
}
