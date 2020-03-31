package tr

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
)

// FromCtx init zero-log and existent span from context.
func FromCtx(ctx context.Context) Interface {
	l := log.Ctx(ctx)

	if span := opentracing.SpanFromContext(ctx); span != nil {
		return NewZeroTracer(l, span)
	}

	return NewZero(l)
}
