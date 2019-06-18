package packager

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/pkg/errors"
)

func (p *Packager) s3BucketName() (string, error) {
	if p.cachedBucketName != "" {
		return p.cachedBucketName, nil
	}

	getAccountResp, err := p.sts.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return "", errors.Wrap(err, "determining aws account id")
	}

	accountId := *getAccountResp.Account
	bucketName := fmt.Sprintf("stackit-%s-%s", p.region, accountId)

	enableVersioning := func() error {
		_, err := p.s3.PutBucketVersioning(&s3.PutBucketVersioningInput{
			Bucket: &bucketName,
			VersioningConfiguration: &s3.VersioningConfiguration{
				Status: aws.String(s3.BucketVersioningStatusEnabled),
			},
		})

		return errors.Wrap(err, "enabling versioning on bucket")
	}

	resp, err := p.s3.GetBucketVersioning(&s3.GetBucketVersioningInput{Bucket: &bucketName})
	if err != nil {
		if err, ok := err.(awserr.Error); ok {
			if err.Code() == s3.ErrCodeNoSuchBucket {
				_, err := p.s3.CreateBucket(&s3.CreateBucketInput{Bucket: &bucketName})
				if err != nil {
					return "", errors.Wrap(err, "creating s3 bucket")
				}

				err = enableVersioning()
				if err != nil {
					return "", err
				}

				p.cachedBucketName = bucketName
				return bucketName, nil
			}
		}
		return "", errors.Wrap(err, "determining if s3 bucket exists")
	}

	if resp.Status == nil || *resp.Status != s3.BucketVersioningStatusEnabled {
		err = enableVersioning()
		if err != nil {
			return "", err
		}
	}

	p.cachedBucketName = bucketName
	return bucketName, nil
}
