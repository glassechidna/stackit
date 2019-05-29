package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnsureJson(t *testing.T) {
	t.Run("input json", func(t *testing.T) {
		output, err := ensureJson(`{"Hello": "world"}`)
		assert.NoError(t, err)
		assert.Equal(t, `{"Hello": "world"}`, string(output))
	})

	t.Run("input yaml", func(t *testing.T) {
		output, err := ensureJson(`Hello: world`)
		assert.NoError(t, err)
		assert.Equal(t, `{"Hello":"world"}`, string(output))
	})
}
