package helpers

import (
	"context"
	"time"
)

var defaultRequestTimeout time.Duration = 15 * time.Second

func CtxDefaultTimeout(ctx context.Context, timeout *time.Duration) (context.Context, context.CancelFunc) {
	if timeout == nil {
		timeout = &defaultRequestTimeout
	}
	return context.WithTimeout(ctx, *timeout)
}
