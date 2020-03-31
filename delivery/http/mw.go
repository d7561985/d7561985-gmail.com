package http

import (
	"fmt"
	"net/http"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/d7561985/questions/internal/tr"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/rs/zerolog"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
)

const fourKB = 4 << 10

func (s *service) Recovery(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				var opt []opentracing.StartSpanOption
				if sp := opentracing.SpanFromContext(r.Context()); sp != nil {
					opt = append(opt, opentracing.ChildOf(sp.Context()))
				}

				span := s.trace.StartSpan("PANIC", opt...)
				ctx := opentracing.ContextWithSpan(r.Context(), span)
				l := tr.FromCtx(ctx)

				defer span.Finish()

				buf := make([]byte, fourKB)
				size := runtime.Stack(buf, true)
				buf = buf[:size]

				l.LogFields(zerolog.ErrorLevel,
					log.String("stacktrace", string(buf)),
					log.String("stack", string(debug.Stack())),
					//log.String("handler", ctx.HandlerName()),
					log.Object("trace", rvr),
				)

				ext.HTTPStatusCode.Set(span, uint16(http.StatusInternalServerError))
				w.WriteHeader(http.StatusInternalServerError)

				fmt.Println(rvr)
				debug.PrintStack()
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (s *service) Jaeger(next http.Handler) http.Handler {
	return nethttp.Middleware(
		s.trace,
		next,
		nethttp.OperationNameFunc(func(r *http.Request) string {
			return "HTTP " + r.Method + r.URL.Path
		}),
		nethttp.MWSpanObserver(func(sp opentracing.Span, r *http.Request) {
			//sp.SetTag("http.uri", r.URL.EscapedPath())
		}),
		nethttp.MWSpanFilter(func(r *http.Request) bool {
			return !(r.Method == http.MethodGet && strings.HasPrefix(r.URL.RequestURI(), "/health"))
		}),
	)
}

func (s *service) Logger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && strings.HasPrefix(r.URL.RequestURI(), "/health") {
			next.ServeHTTP(w, r)
			return
		}

		zl := s.log.BG().With().
			Str("ip", r.RemoteAddr).
			Str("method", r.Method).
			Str("path", r.URL.RequestURI()).
			Str("user-agent", r.Header.Get("User-Agent")).
			Logger()

		_ctx := zl.WithContext(r.Context())
		req := r.WithContext(_ctx)
		next.ServeHTTP(w, req)
	}

	return http.HandlerFunc(fn)
}
