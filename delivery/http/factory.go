package http

import (
	"github.com/d7561985/questions/internal/tr"
	"github.com/d7561985/questions/usecase"
	"github.com/go-chi/chi"
	"github.com/opentracing/opentracing-go"
)

type factory struct {
	log   tr.Interface
	trace opentracing.Tracer
	uc    usecase.Interface
}

func (f *factory) route(v string) chi.Router {
	switch v {
	case questions:
		h := &questionHandler{factory: f}
		return h.route()
	}

	return nil
}
