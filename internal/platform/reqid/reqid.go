package reqid

import (
	"context"
	"crypto/rand"
	"encoding/hex"
)

type ctxKey struct{}

func New() string {
	var b [16]byte
	_, _ = rand.Read(b[:])

	return hex.EncodeToString(b[:])
}

func With(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKey{}, id)
}

func Get(ctx context.Context) string {
	v := ctx.Value(ctxKey{})
	if s, ok := v.(string); ok {
		return s
	}

	return ""
}
