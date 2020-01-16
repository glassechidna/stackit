package packager

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
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

	addTags := func() error {
		tags := getStackitTags()
		_, err := p.s3.PutBucketTagging(&s3.PutBucketTaggingInput{
			Bucket: &bucketName,
			Tagging: &s3.Tagging{
				TagSet: tags,
			},
		})
		return errors.Wrap(err, "adding tags on bucket")
	}

	_, getTagErr := p.s3.GetBucketTagging(&s3.GetBucketTaggingInput{Bucket: &bucketName})
	if tagError, ok := getTagErr.(awserr.Error); ok && tagError.Code() == "NoSuchTagSet" {
		err = addTags()
		if err != nil {
			return "", err
		}
	}

	p.cachedBucketName = bucketName
	return bucketName, nil
}

func getStackitTags() []*s3.Tag {
	tags := viper.GetString("stackit-tags")
	if tags == "" {
		return nil
	}

	tagList := strings.Split(tags,",")
	tagMap := make(map[string]string)
	for _, pair := range tagList {
		dict := strings.Split(pair, "=")
		tagMap[dict[0]] = dict[1]
	}

	result := make([]*s3.Tag, 0, len(tagMap))
	for k, v := range tagMap {
		var t = &s3.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}
		result = append(result, t)
	}

	return result
}
