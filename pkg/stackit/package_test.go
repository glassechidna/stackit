package stackit

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnsureBucket(t *testing.T) {
	t.SkipNow()

	sess := session.Must(session.NewSession(aws.NewConfig().WithRegion("ap-southeast-2")))

	p := &Packager{
		sts:    sts.New(sess),
		s3:     s3.New(sess),
		region: *sess.Config.Region,
	}

	bucket, err := p.s3BucketName()
	assert.NoError(t, err)
	assert.NotEmpty(t, bucket)
	spew.Dump(bucket)
}

func TestExtractTemplateFromCliStdout(t *testing.T) {
	input := []byte(`Uploading to mtdtest/bce72dc454d3a126d60cd2ef6857ce62  2919617 / 2919617.0  (100.00%)
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Parameters:`)

	expected := `AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Parameters:`

	output := extractTemplateFromCliStdout(input)
	assert.Equal(t, expected, output)
}

func TestExtractTemplateFromCliStdout_WithLeadingWhitespace(t *testing.T) {
	input := []byte(`
Uploading to mtdtest/bce72dc454d3a126d60cd2ef6857ce62  2919617 / 2919617.0  (100.00%)
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Parameters:`)

	expected := `AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Parameters:`

	output := extractTemplateFromCliStdout(input)
	assert.Equal(t, expected, output)
}

