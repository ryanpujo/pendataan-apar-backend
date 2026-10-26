package utils

import (
	"context"
	"testing"
)

type TestTable map[string]struct {
	Payload []byte
	URIParam any
	SetupCTX func() (context.Context, context.CancelFunc)
	Arrange func()
	ArrangeContext func(ctx context.Context)
	Assert func(t *testing.T, actual any, err error)
	AssertHTTP func(t *testing.T, code int, response any)
}