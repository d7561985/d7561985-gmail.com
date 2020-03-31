package tr

import (
	"github.com/opentracing/opentracing-go/log"
	"github.com/rs/zerolog"
	"github.com/uber/jaeger-client-go"
)

type Interface interface {
	jaeger.Logger
	BG() *zerolog.Logger
	LogFields(level zerolog.Level, fields ...log.Field)
}
