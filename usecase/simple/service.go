package simple

import (
	"context"
	"fmt"
	"time"

	"github.com/d7561985/questions/internal/tr"
	"github.com/d7561985/questions/model"
	"github.com/d7561985/questions/repository"
	"github.com/d7561985/questions/usecase"
	"github.com/opentracing/opentracing-go/log"
	"github.com/rs/zerolog"
	"golang.org/x/text/language"
)

const (
	TimeFormat = "2006-01-02 15:04:05"
)

type service struct {
	repo       repository.Interface
	translator repository.Translator
	cache      repository.Cache

	log     tr.Interface
	timeout time.Duration
}

func NewService(repo repository.Interface, transl repository.Translator, c repository.Cache,
	log tr.Interface, timeout time.Duration) usecase.Interface {
	return &service{repo: repo, translator: transl, cache: c, log: log, timeout: timeout}
}

func (s service) QuestionList(_ctx context.Context, lang string) (model.QuestionList, error) {
	// validation part should be moved
	switch {
	case lang == "":
		return nil, fmt.Errorf(`"lang" query is required`)
	default:
		// validation: https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes
		if _, err := language.Parse(lang); err != nil {
			return nil, fmt.Errorf("language query: %q not in ISO_639-1 format", lang)
		}
	}

	ctx, cancel := context.WithTimeout(_ctx, s.timeout)
	defer cancel()

	v, err := s.repo.QuestionList(ctx)
	if err != nil {
		return nil, fmt.Errorf("repo retrieve error: %w", err)
	}

	if s.translator == nil {
		s.log.LogFields(zerolog.WarnLevel, log.String("msg", "no translator initiated"))
		return v, nil
	}

	res := make(model.QuestionList, 0, len(v))

	for _, question := range v {
		t, err := s.cache.GetCache(question.Text, lang)
		if err == nil {
			res = append(res, t)
			continue
		}

		t, err = s.translator.TranslateIt(ctx, lang, question)
		if err != nil {
			return nil, fmt.Errorf("question: %q has problem with translation: %w", question.Text, err)
		}

		// SOA oriented on good response. That's why we only notify about problem on possible cache side
		// without interruption good response.
		if err := s.cache.Put(question.Text, lang, t); err != nil {
			s.log.LogFields(zerolog.ErrorLevel, log.Error(err),
				log.String("question", question.Text), log.String("lang", lang))
		}

		res = append(res, t)
	}

	return res, nil
}

func (s service) AddQuestion(_ctx context.Context, question model.Question) (model.Question, error) {
	ctx, cancel := context.WithTimeout(_ctx, s.timeout)
	defer cancel()

	question.CreatedAt = time.Now().Format(TimeFormat)
	if err := s.repo.AddQuestion(ctx, question); err != nil {
		return question, fmt.Errorf("save question %q error: %w", question.Text, err)
	}

	return question, nil
}

// should check if repository is OK
func (s service) Health(ctx context.Context) error {
	return nil
}
