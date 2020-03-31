package tr

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/rs/zerolog"
)

type ZT struct {
	*Z
	span opentracing.Span
}

func NewZeroTracer(l *zerolog.Logger, span opentracing.Span) Interface {
	return &ZT{Z: NewZero(l), span: span}
}

// LogFields is an efficient and type-checked way to record key:value
// logging data about a Span, though the programming interface is a little
// more verbose than LogKV(). Here's an example:
//
//    span.LogFields(
//        log.String("event", "soft error"),
//        log.String("type", "cache timeout"),
//        log.Int("waited.millis", 1500))
//
// Also see Span.FinishWithOptions() and FinishOptions.BulkLogData.
func (z *ZT) LogFields(level zerolog.Level, fields ...log.Field) {
	z.Z.LogFields(level, fields...)

	if len(fields) == 0 {
		return
	}

	z.span.LogFields(fields...)
	// it's redundant to have 2 ranges.
	for _, f := range fields {
		if _, ok := f.Value().(error); ok {
			ext.Error.Set(z.span, true)
		}
	}
}
