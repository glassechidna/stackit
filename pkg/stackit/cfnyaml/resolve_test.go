package cfnyaml

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestResolve(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/aliases_input.yml")
	assert.NoError(t, err)

	c, err := Parse(b)
	assert.NoError(t, err)

	b, err = ioutil.ReadFile("testdata/aliases_expected.yml")
	assert.NoError(t, err)
	expected := string(b)

	assert.Equal(t, expected, c.String())
}

