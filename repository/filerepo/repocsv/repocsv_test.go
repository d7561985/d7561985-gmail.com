// +build unit

package repocsv

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/d7561985/questions/model"
	"github.com/stretchr/testify/assert"
)

func Test_service_Read(t *testing.T) {
	f, err := os.OpenFile("questions.csv", os.O_RDWR, os.ModePerm)
	assert.NoError(t, err)

	x := make(model.QuestionList, 0)
	srv := New(f)
	res, err := srv.Read(x)
	assert.NoError(t, err)

	fmt.Println(res)

	x = res.(model.QuestionList)
	x = append(x, model.Question{
		Text:      "test",
		CreatedAt: time.Now(),
		Choices:   [3]model.Choice{{"1"}, {"2"}, {"3"}},
	})

	err = srv.Write(x)
	assert.NoError(t, err)
}
