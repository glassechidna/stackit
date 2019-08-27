package packager

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

func TestMd5path(t *testing.T) {
	sum, err := md5path("testdata/helloworld.txt")
	assert.NoError(t, err)
	assert.Equal(t, "6f5902ac237024bdd0c176cb93063dc4", sum)
}
