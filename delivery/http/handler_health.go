package http

import (
	"net/http"

	"github.com/d7561985/questions/internal/tr"
	"github.com/opentracing/opentracing-go/log"
	"github.com/prometheus/common/version"
	"github.com/rs/zerolog"
)

func (s *service) HealthHandler(w http.ResponseWriter, r *http.Request) {
	if err := s.uc.Health(r.Context()); err != nil {
		tr.FromCtx(r.Context()).LogFields(zerolog.ErrorLevel, log.Error(err), log.String("event", "health"))
		Send(w, r.Header.Get("Content-Type"), nil, http.StatusBadGateway, s.log, err)

		return
	}

	Send(w, r.Header.Get("Content-Type"), map[string]interface{}{"build": version.BuildContext(), "version": version.Info()}, http.StatusOK, s.log, nil)
}
