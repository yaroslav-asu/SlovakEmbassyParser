package funcs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStripString(t *testing.T) {
	s := "\n\n\t  test \n \t\t"
	stripedString := StripString(s)
	assert.Equal(t, "test", stripedString)

	s = "test test"
	stripedString = StripString(s)
	assert.Equal(t, "test test", stripedString)

	s = "test test\n"
	stripedString = StripString(s)
	assert.Equal(t, "test test", stripedString)
}

func TestStripStringRunes(t *testing.T) {
	s := "test\n"
	stripedString := stripStringRunes(s, '\n')
	assert.Equal(t, "test", stripedString)
	stripedString = stripStringRunes(s, '\t')
	assert.Equal(t, "test\n", stripedString)

	s = "\ntest\n \n"
	stripedString = stripStringRunes(s, '\n', ' ')
	assert.Equal(t, "test", stripedString)
}
