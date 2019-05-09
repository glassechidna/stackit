package cmd

import (
	"bytes"
	"github.com/glassechidna/stackit/pkg/stackit"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"regexp"
	"testing"
)

func TestChangeSetFormatting(t *testing.T) {
	ymlBody := `
input:
  capabilities:
  - CAPABILITY_IAM
  - CAPABILITY_NAMED_IAM
  changesetname: aidan-mtd-test-csid-1557355052
  changesettype: CREATE
  clienttoken: stackit-1586111999
  description: null
  notificationarns: []
  parameters:
  - parameterkey: VatBetaOptinUIUrl
    parametervalue: def
    resolvedvalue: null
    usepreviousvalue: null
  resourcetypes: []
  rolearn: null
  rollbackconfiguration: null
  stackname: aidan-mtd-test
  tags:
  - key: moo
    value: koo
  templatebody: |2
    AWSTemplateFormatVersion: '2010-09-09'
  templateurl: null
  useprevioustemplate: false
output:
  id: arn:aws:cloudformation:ap-southeast-2:720884384464:changeSet/aidan-mtd-test-csid-1557355052/dc7928df-d27e-4992-a350-9ee4ba357999
  stackid: arn:aws:cloudformation:ap-southeast-2:720884384464:stack/aidan-mtd-test/de4820c0-71e1-11e9-9ce6-0230e8e86e24
changes:
- resourcechange:
    action: Add
    details: []
    logicalresourceid: Cell
    physicalresourceid: null
    replacement: null
    resourcetype: Custom::XeroCellInfo
    scope: []
  type: Resource
- resourcechange:
    action: Add
    details: []
    logicalresourceid: CodeDeployServiceRole
    physicalresourceid: null
    replacement: null
    resourcetype: AWS::IAM::Role
    scope: []
  type: Resource
templatebody: ""
`

	input := stackit.PrepareOutput{}
	err := yaml.Unmarshal([]byte(ymlBody), &input)
	assert.NoError(t, err)

	expected := `
Stack ID: arn:aws:cloudformation:ap-southeast-2:720884384464:stack/aidan-mtd-test/de4820c0-71e1-11e9-9ce6-0230e8e86e24
Change Set ID: arn:aws:cloudformation:ap-southeast-2:720884384464:changeSet/aidan-mtd-test-csid-1557355052/dc7928df-d27e-4992-a350-9ee4ba357999
Changes:

\+--------\+-----------------------\+
| ACTION |       RESOURCE        |
\+--------\+-----------------------\+
| Add    | Cell                  |
| Add    | CodeDeployServiceRole |
\+--------\+-----------------------\+
`
	assert.Equal(t, expected, userFriendlyChangesOutput(&input))
}

func TestPackageAndExecuteE2E(t *testing.T) {
	buf := &bytes.Buffer{}
	out := io.MultiWriter(buf, os.Stderr)
	RootCmd.SetOutput(out)

	t.Run("package", func(t *testing.T) {
		RootCmd.SetArgs([]string{
			"package",
			"--stack-name", "test-stack-packaged",
			"--template", "../sample/serverless.yml",
		})
		_ = RootCmd.Execute()

		assert.Regexp(t, regexp.MustCompile(`
Stack ID: arn:aws:cloudformation:ap-southeast-2:\d+:stack/test-stack-packaged/[a-f0-9-]+
Change Set ID: arn:aws:cloudformation:ap-southeast-2:\d+:changeSet/test-stack-packaged-csid-\d+/[a-f0-9-]+
Changes:

\+--------\+---------------------------\+
\| ACTION \|         RESOURCE          \|
\+--------\+---------------------------\+
\| Add    \| Function                  \|
\| Add    \| FunctionAliaslive         \|
\| Add    \| FunctionRole              \|
\| Add    \| FunctionVersion\S+\s*\|
\+--------\+---------------------------\+
`), buf.String())
	})

	buf.Reset()

	t.Run("execute", func(t *testing.T) {
		RootCmd.SetArgs([]string{
			"execute",
			"--stack-name", "test-stack-packaged",
		})
		_ = RootCmd.Execute()

		assert.Regexp(t, regexp.MustCompile(`\[\d\d:\d\d:\d\d]  test-stack-packaged - CREATE_IN_PROGRESS - User Initiated
\[\d\d:\d\d:\d\d]         FunctionRole - CREATE_IN_PROGRESS 
\[\d\d:\d\d:\d\d]         FunctionRole - CREATE_IN_PROGRESS - Resource creation Initiated
\[\d\d:\d\d:\d\d]         FunctionRole - CREATE_COMPLETE 
\[\d\d:\d\d:\d\d]             Function - CREATE_IN_PROGRESS 
\[\d\d:\d\d:\d\d]             Function - CREATE_IN_PROGRESS - Resource creation Initiated
\[\d\d:\d\d:\d\d]             Function - CREATE_COMPLETE 
\[\d\d:\d\d:\d\d] FunctionVersion\S+ - CREATE_IN_PROGRESS 
\[\d\d:\d\d:\d\d] FunctionVersion\S+ - CREATE_IN_PROGRESS - Resource creation Initiated
\[\d\d:\d\d:\d\d] FunctionVersion\S+ - CREATE_COMPLETE 
\[\d\d:\d\d:\d\d]    FunctionAliaslive - CREATE_IN_PROGRESS 
\[\d\d:\d\d:\d\d]    FunctionAliaslive - CREATE_IN_PROGRESS - Resource creation Initiated
\[\d\d:\d\d:\d\d]    FunctionAliaslive - CREATE_COMPLETE`), buf.String())
	})

	buf.Reset()

	t.Run("down", func(t *testing.T) {
		RootCmd.SetArgs([]string{
			"down",
			"--stack-name", "test-stack-packaged",
		})
		_ = RootCmd.Execute()
	})
}
