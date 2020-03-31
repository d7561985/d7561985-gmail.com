package filerepo

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/d7561985/questions/model"
	"github.com/jszwec/csvutil"
	"github.com/stretchr/testify/assert"
)

func TestX(t *testing.T) {
	x := [][]string{
		{"Question text", "Created At", "Choice 1", "Choice", "Choice 3"},
		{"What is the capital of Luxembourg ?", "2019-06-01 00:00:00", "Luxembourg", "Paris", "Berlin"},
		{"What does mean O.A.T. ?", "2019-06-02 00:00:00", "Open Assignment Technologies", "Open Assessment Technologies", "Open Acknowledgment Technologies"},
	}

	f, err := os.OpenFile("x.csv", os.O_CREATE|os.O_WRONLY, os.ModePerm)
	assert.NoError(t, err)

	w := csv.NewWriter(f)
	err = w.WriteAll(x)
	assert.NoError(t, err)
}

func TestAA(t *testing.T) {
	f, err := ioutil.ReadFile("questions.csv")
	assert.NoError(t, err)

	reader := bytes.NewBuffer(f)
	r := csv.NewReader(reader)
	r.LazyQuotes = true

	/*for {
		v, err := r.Read()
		if err != nil {
			assert.NoError(t, err)
			break
		}

		fmt.Println(v)
	}*/

	d, err := csvutil.NewDecoder(r)
	assert.NoError(t, err)
	d.Map = func(field, col string, v interface{}) string {
		fmt.Println(field, col, v)
		return field
	}

	s := make(model.QuestionList, 0)

	err = d.Decode(&s)
	assert.NoError(t, err)
	fmt.Println(s)
}
