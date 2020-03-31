// +build unit

package simple

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestLanguage(t *testing.T) {
	list := []string{"ru", "uk", "en", "be"}

	for _, lang := range list {
		_, err := language.Parse(lang)
		assert.NoError(t, err)
	}
}
