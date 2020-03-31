package translate

import (
	"context"
	"fmt"
	"time"

	translate "cloud.google.com/go/translate/apiv3"
	"github.com/d7561985/questions/model"
	"github.com/d7561985/questions/repository"
	"github.com/opentracing/opentracing-go"
	translatepb "google.golang.org/genproto/googleapis/cloud/translate/v3"
)

type service struct {
	trace     opentracing.Tracer
	client    *translate.TranslationClient
	projectID string
}

func New(pID string, trace opentracing.Tracer, timeout time.Duration) (repository.Translator, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	c, err := translate.NewTranslationClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("create translate client: %w", err)
	}

	return &service{projectID: pID, trace: trace, client: c}, nil
}

func (s *service) TranslateIt(ctx context.Context, lang string, question model.Question) (model.Question, error) {
	end := s.start(ctx, "translate")
	defer end()

	req := &translatepb.TranslateTextRequest{
		Parent:             fmt.Sprintf("projects/%s/locations/global", s.projectID),
		Contents:           []string{question.Text, question.Choices[0].Text, question.Choices[1].Text, question.Choices[2].Text},
		MimeType:           "text/plain", // Mime types: "text/plain", "text/html"
		SourceLanguageCode: "en-US",
		TargetLanguageCode: lang,
		Model:              "",
		GlossaryConfig:     nil,
		Labels:             nil,
	}

	resp, err := s.client.TranslateText(ctx, req)
	if err != nil {
		return question, fmt.Errorf("TranslateText: %v", err)
	}

	// Display the translation for each input text provided
	for i, translation := range resp.GetTranslations() {
		t := translation.GetTranslatedText()

		switch i {
		case 0:
			question.Text = t
		case 1:
			question.Choices[0].Text = t
		case 2:
			question.Choices[1].Text = t
		case 3:
			question.Choices[2].Text = t
		}

		// fmt.Println(i, t)
	}

	return question, nil
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
