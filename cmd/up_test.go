package cmd

import (
	"bytes"
	"github.com/magiconair/properties/assert"
	"io"
	"os"
	"strings"
	"testing"
)

func TestUp(t *testing.T) {
	if testing.Short() {
		t.Skip("skip e2e tests in short mode")
	}

	t.Run("up", func(t *testing.T) {
		RootCmd.SetArgs([]string{
			"up",
			"--stack-name", "test-stack",
			"--template", "../sample/sample.yml",
		})

		buf := &bytes.Buffer{}
		out := io.MultiWriter(buf, os.Stderr)
		RootCmd.SetOutput(out)

		_ = RootCmd.Execute()

		assert.Matches(t, buf.String(), strings.TrimSpace(`
\[\d\d:\d\d:\d\d]           test-stack - CREATE_IN_PROGRESS - User Initiated
\[\d\d:\d\d:\d\d]             LogGroup - CREATE_IN_PROGRESS 
\[\d\d:\d\d:\d\d]             LogGroup - CREATE_IN_PROGRESS - Resource creation Initiated
\[\d\d:\d\d:\d\d]             LogGroup - CREATE_COMPLETE 
\[\d\d:\d\d:\d\d]              TaskDef - CREATE_IN_PROGRESS 
\[\d\d:\d\d:\d\d]              TaskDef - CREATE_IN_PROGRESS - Resource creation Initiated
\[\d\d:\d\d:\d\d]              TaskDef - CREATE_COMPLETE 
\[\d\d:\d\d:\d\d]          TargetGroup - CREATE_IN_PROGRESS 
\[\d\d:\d\d:\d\d]          TargetGroup - CREATE_IN_PROGRESS - Resource creation Initiated
\[\d\d:\d\d:\d\d]          TargetGroup - CREATE_COMPLETE 
\{
  "LogGroup": "test-stack-LogGroup",
  "TaskDef": "arn:aws:ecs:ap-southeast-2:607481581596:task-definition/ecs-run-task-test:\d+"
\}
`))
	})

	t.Run("down", func(t *testing.T) {
		RootCmd.SetArgs([]string{
			"down",
			"--stack-name", "test-stack",
		})

		buf := &bytes.Buffer{}
		out := io.MultiWriter(buf, os.Stderr)
		RootCmd.SetOutput(out)

		_ = RootCmd.Execute()

		assert.Matches(t, buf.String(), strings.TrimSpace(`
\[\d\d:\d\d:\d\d]           test-stack - DELETE_IN_PROGRESS - User Initiated
\[\d\d:\d\d:\d\d]          TargetGroup - DELETE_IN_PROGRESS 
\[\d\d:\d\d:\d\d]          TargetGroup - DELETE_COMPLETE 
\[\d\d:\d\d:\d\d]              TaskDef - DELETE_IN_PROGRESS 
\[\d\d:\d\d:\d\d]              TaskDef - DELETE_COMPLETE 
\[\d\d:\d\d:\d\d]             LogGroup - DELETE_IN_PROGRESS 
\[\d\d:\d\d:\d\d]             LogGroup - DELETE_COMPLETE
`))
	})

}
