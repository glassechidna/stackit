package packager

import (
	"context"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

type UploadedObject struct {
	Bucket        string
	Key           string
	VersionId     string
	AlreadyExists bool
}

func (p *Packager) Upload(ctx context.Context, key, path string) (*UploadedObject, error) {
	uploader := s3manager.NewUploaderWithClient(p.s3)

	bucket, err := p.s3BucketName()
	if err != nil {
		return nil, errors.Wrap(err, "getting bucket for object upload")
	}

	head, err := p.s3.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err == nil {
		ver := ""
		if head.VersionId != nil {
			ver = *head.VersionId
		}

		return &UploadedObject{
			Bucket:        bucket,
			Key:           key,
			VersionId:     ver,
			AlreadyExists: true,
		}, nil
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, errors.Wrapf(err, "determining absolute path of '%s'", path)
	}

	file, err := os.Open(absPath)
	if err != nil {
		return nil, errors.Wrapf(err, "opening file '%s'", absPath)
	}

	resp, err := uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   file,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "uploading %s to to s3://%s/%s", absPath, bucket, key)
	}

	return &UploadedObject{
		Bucket:    bucket,
		Key:       key,
		VersionId: *resp.VersionID,
	}, nil
}

//func (p *Packager) Head(ctx context.Context, key string) bool {
//	bucket, err := p.s3BucketName()
//	if err != nil {
//		return false
//	}
//
//	_, err = p.s3.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
//		Bucket: &bucket,
//		Key:    &key,
//	})
//	if err != nil {
//		return false
//	}
//
//	return true
//}
