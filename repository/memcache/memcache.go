package memcache

import (
	"errors"

	"github.com/d7561985/questions/model"
	"github.com/d7561985/questions/repository"
)

var errNoInCache = errors.New("not found")

// crude memcache imp
type service struct {
	store map[string]model.Question
}

func New() repository.Cache {
	return &service{store: make(map[string]model.Question)}
}

func (s *service) makeKey(questionID, lang string) string {
	return questionID + lang
}

func (s *service) GetCache(questionID, lang string) (model.Question, error) {
	res, ok := s.store[s.makeKey(questionID, lang)]
	if !ok {
		return res, errNoInCache
	}

	return res, nil
}

func (s *service) Put(questionID, lang string, in model.Question) error {
	s.store[s.makeKey(questionID, lang)] = in
	return nil
}
