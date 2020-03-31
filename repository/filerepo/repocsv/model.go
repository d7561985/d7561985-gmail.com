package repocsv

import (
	"github.com/d7561985/questions/model"
)

type QuestionList []Question

type Question struct {
	Test      string `csv:"Question text"`
	CreatedAt string `csv:"Created At"`
	Choice1   string `csv:"Choice 1"`
	Choice2   string `csv:"Choice"`
	Choice3   string `csv:"Choice 3"`
}

func NewQuestionFromModel(in model.Question) Question {
	return Question{
		Test:      in.Text,
		CreatedAt: in.CreatedAt,
		Choice1:   in.Choices[0].Text,
		Choice2:   in.Choices[1].Text,
		Choice3:   in.Choices[2].Text,
	}
}

func (q *Question) Model() model.Question {
	return model.Question{
		Text:      q.Test,
		CreatedAt: q.CreatedAt,
		Choices:   [3]model.Choice{{Text: q.Choice1}, {Text: q.Choice2}, {Text: q.Choice3}},
	}
}

func NewQuestionListFromModel(in model.QuestionList) QuestionList {
	res := make(QuestionList, len(in))

	for i := range in {
		res[i] = NewQuestionFromModel(in[i])
	}

	return res
}

func (q QuestionList) Model() model.QuestionList {
	res := make(model.QuestionList, len(q))

	for i, question := range q {
		res[i] = question.Model()
	}

	return res
}
