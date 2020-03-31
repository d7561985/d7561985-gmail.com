package filerepo

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/d7561985/questions/internal/tr"
	"github.com/d7561985/questions/model"
	"github.com/d7561985/questions/repository"
	"github.com/d7561985/questions/repository/filerepo/repocsv"
	"github.com/d7561985/questions/repository/filerepo/repojson"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/rs/zerolog"
)

const folder = "data"

type service struct {
	trace opentracing.Tracer
	log   tr.Interface
	store model.QuestionList

	file repository.FileInterface
}

func New(log tr.Interface, trace opentracing.Tracer) *service {
	return &service{log: log, trace: trace, store: make(model.QuestionList, 0)}
}

func (s *service) QuestionList(ctx context.Context) (model.QuestionList, error) {
	closer := s.start(ctx, "Read repo")
	defer closer()

	return s.store, nil
}

func (s *service) AddQuestion(ctx context.Context, question model.Question) error {
	done := make(chan error)

	go func() {
		s.store = append(s.store, question)
		if err := s.file.Write(s.store); err != nil {
			done <- err
			return
		}

		done <- nil
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

func (s *service) Load(file string) error {
	p := path.Join(folder, file)

	f, err := os.OpenFile(p, os.O_RDWR, os.ModePerm)
	if err != nil {
		return fmt.Errorf("open file %q error: %w", p, err)
	}

	switch ext := path.Ext(file); ext {
	case ".json":
		s.file = repojson.New(f)
	case ".csv":
		s.file = repocsv.New(f)
	default:
		return fmt.Errorf("unsupported extension %s file %q", ext, p)
	}

	res, err := s.file.Read(&s.store)
	if err != nil {
		return fmt.Errorf("encode file %s error: %w", p, err)
	}

	// we trust in cast
	s.store = res.(model.QuestionList)

	s.log.LogFields(zerolog.InfoLevel, log.String("load file", p))

	return nil
}

func (s *service) Close() error {
	if s.file == nil {
		return nil
	}

	return s.file.Close()
}

// graceful open-tracing log system
func (s *service) start(_ctx context.Context, operation string) func() {
	opt := make([]opentracing.StartSpanOption, 0, 1)

	span := opentracing.SpanFromContext(_ctx)
	if span != nil {
		opt = append(opt, opentracing.ChildOf(span.Context()))
	}

	span = s.trace.StartSpan(operation, opt...)

	return func() {
		span.Finish()
	}
}
