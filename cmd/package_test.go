package cmd

import (
	"bytes"
	"github.com/glassechidna/stackit/pkg/stackit"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"testing"
	"time"
)

func TestUsefulErrorIfTemplateDoesntExist(t *testing.T) {
	buf := &bytes.Buffer{}
	RootCmd.SetOutput(buf)

	RootCmd.SetArgs([]string{
		"package",
		"--stack-name", "some-stack-name",
		"--template", "doesnt-exist.yml",
	})

	assert.NotPanics(t, func() {
		RootCmd.Execute()
	})

	assert.Regexp(t, regexp.MustCompile(`^no file exists at`), buf.String())
}

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

+--------+-----------------------+----------------------+
| ACTION |       RESOURCE        |         TYPE         |
+--------+-----------------------+----------------------+
| Add    | Cell                  | Custom::XeroCellInfo |
| Add    | CodeDeployServiceRole | AWS::IAM::Role       |
+--------+-----------------------+----------------------+
`
	assert.Equal(t, expected, userFriendlyChangesOutput(&input))
}

func TestPackageAndExecuteE2E(t *testing.T) {
	if testing.Short() {
		t.Skip("skip e2e tests in short mode")
	}

	buf := &bytes.Buffer{}
	out := io.MultiWriter(buf, os.Stderr)
	RootCmd.SetOutput(out)

	t.Run("up", func(t *testing.T) {
		// to ensure the first run doesn't report "object already exists"
		err := ioutil.WriteFile("../sample/func/random.txt", []byte(time.Now().String()), 0644)
		assert.NoError(t, err)

		RootCmd.SetArgs([]string{
			"up",
			"--stack-name", "test-stack-packaged",
			"--template", "../sample/serverless.yml",
		})
		_ = RootCmd.Execute()

		assert.Regexp(t, regexp.MustCompile(`Uploaded ./func to s3://stackit-ap-southeast-2-607481581596/test-stack-packaged/func.zip/[a-f0-9]{32} \(v = [^)]+\)
\[\d\d:\d\d:\d\d]  test-stack-packaged - CREATE_IN_PROGRESS - User Initiated
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
\[\d\d:\d\d:\d\d]    FunctionAliaslive - CREATE_COMPLETE 
\[\d\d:\d\d:\d\d]  test-stack-packaged - CREATE_COMPLETE 
`), buf.String())
	})

	buf.Reset()

	t.Run("up second time", func(t *testing.T) {
		RootCmd.SetArgs([]string{
			"up",
			"--stack-name", "test-stack-packaged",
			"--template", "../sample/serverless.yml",
		})
		_ = RootCmd.Execute()

		assert.Regexp(t, regexp.MustCompile(`./func already exists at s3://stackit-ap-southeast-2-607481581596/test-stack-packaged/func.zip/[a-f0-9]{32} \(v = [^)]+\)`), buf.String())
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
