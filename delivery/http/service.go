package http

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/d7561985/questions/delivery"
	"github.com/d7561985/questions/internal/tr"
	"github.com/d7561985/questions/usecase"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/rs/zerolog"
)

const (
	timeOut = 5
)

type service struct {
	server *http.Server
	*factory
}

func New(log tr.Interface, uc usecase.Interface, tracer opentracing.Tracer) delivery.Interface {
	return &service{
		server: &http.Server{},
		factory: &factory{
			log:   log,
			uc:    uc,
			trace: tracer,
		},
	}
}

func (s *service) Serve(listener net.Listener) {
	s.server.Handler = s.Router()

	s.log.LogFields(zerolog.InfoLevel, log.String("addr", listener.Addr().String()),
		log.String("event", "http delivery start"))

	if err := s.server.Serve(listener); err != nil {
		s.log.LogFields(zerolog.ErrorLevel, log.Error(err))
	}
}

func (s *service) Stop() {
	if s.server == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*timeOut)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.log.LogFields(zerolog.ErrorLevel, log.Error(err))
	}
}

// send RESTfull convention with requested context type in header field
// in default Content-Type cases handle json
func Send(w http.ResponseWriter, ct string, req interface{}, status int, l tr.Interface, er error) {
	// fast error sender
	if er != nil {
		req = map[string]string{"error": er.Error()}
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(status)
	}

	switch ct {
	case "application/xml":
		w.Header().Add("Content-Type", ct)

		e := xml.NewEncoder(w)
		if err := e.Encode(req); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			l.LogFields(zerolog.ErrorLevel, log.String("event", "marshal"), log.Error(err))

			return
		}
	case "application/json":
		fallthrough
	default:
		w.Header().Add("Content-Type", "application/json")

		e := json.NewEncoder(w)
		if err := e.Encode(req); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			l.LogFields(zerolog.ErrorLevel, log.String("event", "marshal"), log.Error(err))

			return
		}
	}

	if er != nil {
		l.LogFields(zerolog.InfoLevel, log.String("error", er.Error()))
	}
}

// Read post data according with Content-Type provided in header
// by default system read json
func Read(r *http.Request, model interface{}) error {
	switch ct := r.Header.Get("Content-Type"); ct {
	case "application/xml":
		d := xml.NewDecoder(r.Body)
		if err := d.Decode(model); err != nil {
			return fmt.Errorf("unmarshal xml body error: %w", err)
		}
	case "application/json":
		fallthrough
	default:
		d := json.NewDecoder(r.Body)
		if err := d.Decode(model); err != nil {
			return fmt.Errorf("unmarshal json body error: %w", err)
		}
	}

	return nil
}
