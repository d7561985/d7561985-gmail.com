package model

const ChoiceCount = 3

type Choice struct {
	Text string `json:"text"`
}

type Question struct {
	Text      string              `json:"text"`
	CreatedAt string              `json:"createdAt" `
	Choices   [ChoiceCount]Choice `json:"choices"`
}

type QuestionList []Question
