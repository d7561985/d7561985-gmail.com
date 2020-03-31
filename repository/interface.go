//go:generate mockery  -all

//Package repository represent persistent layer with ease scalability to different repos.
package repository

import (
	"context"

	"github.com/d7561985/questions/model"
)

type Interface interface {
	QuestionList(ctx context.Context) (model.QuestionList, error)
	AddQuestion(ctx context.Context, question model.Question) error
	Close() error
}

type Translator interface {
	TranslateIt(ctx context.Context, lang string, uestion model.Question) (model.Question, error)
}

// For preventing translate API from overuse we should use cache.
// cache can be represent as well in memory object as redis implementation
type Cache interface {
	// @questionID - we use quentin text here for simplicity
	GetCache(questionID, lang string) (model.Question, error)
	Put(questionID, lang string, in model.Question) error
}

// Internal generic interface for file repository for handle and support any desired format
type FileInterface interface {
	// Require pointer for *model.QuestionList
	Read(storage interface{}) (interface{}, error)

	// No pointer for model.QuestionList
	Write(storage interface{}) error

	Close() error
}
