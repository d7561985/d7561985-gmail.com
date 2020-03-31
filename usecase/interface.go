//go:generate mockery  -all
package usecase

import (
	"context"

	"github.com/d7561985/questions/model"
)

type Interface interface {
	Health(ctx context.Context) error
	QuestionList(ctx context.Context, lang string) (model.QuestionList, error)
	AddQuestion(ctx context.Context, question model.Question) (model.Question, error)
}
