package repojson

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/d7561985/questions/model"
	"github.com/d7561985/questions/repository"
)

type service struct {
	file io.ReadWriteCloser
}

func New(file *os.File) repository.FileInterface {
	return &service{file: file}
}

func (j *service) Read(storage interface{}) (interface{}, error) {
	e := json.NewDecoder(j.file)
	if err := e.Decode(storage); err != nil {
		return nil, fmt.Errorf("encode error: %w", err)
	}

	switch in := storage.(type) {
	case *model.QuestionList:
		return *in, nil
	default:
		panic(fmt.Errorf("model %T not supported", in))
	}
}

func (j *service) Write(storage interface{}) error {
	w := json.NewEncoder(j.file)
	if err := w.Encode(storage); err != nil {
		return fmt.Errorf("decode error: %w", err)
	}

	return nil
}

func (j *service) Close() error {
	return j.file.Close()
}
