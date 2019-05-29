package cfnyaml

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

type rewrittenLocation struct {
	bucket    string
	key       string
	versionId string
}

type tableTestEntry struct {
	Name         string
	Explanation  string
	Replacements map[string]rewrittenLocation
}

var tests = []tableTestEntry{
	{
		Name:        "a",
		Explanation: "aws::serverless::function",
		Replacements: map[string]rewrittenLocation{
			"./func": {"bucket", "key.zip", "version"},
		},
	},
	{
		Name:        "b",
		Explanation: "aws::lambda::function",
		Replacements: map[string]rewrittenLocation{
			"./func": {"bucket", "key.zip", "version"},
		},
	},
	{
		Name:        "c",
		Explanation: "aws::lambda::function without version id",
		Replacements: map[string]rewrittenLocation{
			"./func": {"bucket", "key.zip", ""},
		},
	},
	{
		Name:        "d",
		Explanation: "aws::cloudformation::stack without version id",
		Replacements: map[string]rewrittenLocation{
			"./stack.yml": {"bucket", "key.yml", ""},
		},
	},
	{
		Name:        "e",
		Explanation: "aws::cloudformation::stack",
		Replacements: map[string]rewrittenLocation{
			"./stack.yml": {"bucket", "key.yml", "abc"},
		},
	},
}

func TestCfnYaml_PackageableNodes(t *testing.T) {
	for _, e := range tests {
		t.Run(e.Explanation, func(t *testing.T) {
			c, err := ParseFile(fmt.Sprintf("testdata/%s_input.yml", e.Name))
			assert.NoError(t, err)

			nodes, err := c.PackageableNodes()
			assert.NoError(t, err)

			for _, n := range nodes {
				replacement, found := e.Replacements[n.Value]
				assert.True(t, found)
				n.Replace(replacement.bucket, replacement.key, replacement.versionId)
			}

			expected, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s_expected.yml", e.Name))
			assert.NoError(t, err)
			assert.Equal(t, string(expected), c.String())
		})
	}
}
