package tr

import (
	"os"

	"github.com/opentracing/opentracing-go/log"
	"github.com/rs/zerolog"
)

type Z struct {
	log *zerolog.Logger
}

func NewZero(l *zerolog.Logger) *Z {
	return &Z{log: l}
}

func (z *Z) BG() *zerolog.Logger {
	return z.log
}

func (z *Z) Error(msg string) {
	z.log.Error().Msg(msg)
}

func (z *Z) Infof(msg string, args ...interface{}) {
	z.log.Info().Msgf(msg, args...)
}

func (z *Z) LogFields(level zerolog.Level, fields ...log.Field) {
	l := z.log.WithLevel(level)

	for _, f := range fields {
		switch value := f.Value().(type) {
		case string:
			l = l.Str(f.Key(), value)
		case int:
			l = l.Int(f.Key(), value)
		case bool:
			l = l.Bool(f.Key(), value)
		case float64:
			l = l.Float64(f.Key(), value)
		case error:
			l = l.Err(value)
		case interface{}:
			l = l.Interface(f.Key(), value)
		default:
			z.log.Panic().Msgf("LogFields type:%T not supported", value)
		}
	}

	l.Msg("")

	// because of zerolog.WithLevel not handle panic or fatal level we should do it.
	switch level {
	case zerolog.FatalLevel:
		os.Exit(1)
	case zerolog.PanicLevel:
		panic(false)
	}
}
