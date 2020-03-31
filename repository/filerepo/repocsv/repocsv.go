package repocsv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	csvtag "github.com/artonge/go-csv-tag/v2"
	"github.com/d7561985/questions/model"
	"github.com/d7561985/questions/repository"
)

type service struct {
	file io.ReadWriteCloser
}

func New(file *os.File) repository.FileInterface {
	return &service{file: file}
}

// Please ask me why i'm so an angry?
func (s service) Read(storage interface{}) (interface{}, error) {
	r := csv.NewReader(s.file)
	r.LazyQuotes = true

	switch in := storage.(type) {
	case *model.QuestionList:
		x := make(QuestionList, 0)

		if err := csvtag.LoadFromReader(s.file, &x); err != nil {
			return nil, fmt.Errorf("read error %w", err)
		}

		return x.Model(), nil
	default:
		panic(fmt.Errorf("model %T not supported", in))

	}
}

func (s service) Write(storage interface{}) error {
	switch d := storage.(type) {
	case model.QuestionList:
		m := NewQuestionListFromModel(d)
		if err := csvtag.DumpToWriter(m, s.file); err != nil {
			return fmt.Errorf("write errro: %w", err)
		}

	default:
		panic(fmt.Errorf("model %T not supported", d))
	}

	return nil
}

func (s service) Close() error {
	if s.file == nil {
		return nil
	}

	return s.file.Close()
}
