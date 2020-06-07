package cmd

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOutputs(t *testing.T) {
	if testing.Short() {
		t.Skip("skip e2e tests in short mode")
	}

	buf := &bytes.Buffer{}

	t.Run("up", func(t *testing.T) {
		RootCmd.SetArgs([]string{
			"up",
			"--stack-name", "test-stack-s3",
			"--template", "../sample/s3.yml",
		})
		_ = RootCmd.Execute()
	})

	buf.Reset()

	t.Run("outputs", func(t *testing.T) {
		RootCmd.SetArgs([]string{
			"outputs",
			"--stack-name", "test-stack-s3",
		})

		buf := &bytes.Buffer{}
		out := io.MultiWriter(buf, os.Stderr)
		RootCmd.SetOutput(out)

		_ = RootCmd.Execute()

		assert.Regexp(t, regexp.MustCompile(`\{
  "S3BucketName": "test-stack-s3-s3bucket-[a-z0-9]*"
\}
`), buf.String())
	})

	t.Run("down", func(t *testing.T) {
		RootCmd.SetArgs([]string{
			"down",
			"--stack-name", "test-stack-s3",
		})
		_ = RootCmd.Execute()
	})
}
