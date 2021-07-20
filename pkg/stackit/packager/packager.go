package packager

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/glassechidna/stackit/pkg/stackit/cfnyaml"
	"github.com/glassechidna/stackit/pkg/zipper"
	"github.com/pkg/errors"
)

type Packager struct {
	s3       s3iface.S3API
	sts      stsiface.STSAPI
	region   string
	s3Suffix string
	s3Tags   string

	cachedBucketName string
}

func New(s3 s3iface.S3API, sts stsiface.STSAPI, region string, s3Suffix string, s3Tags string) *Packager {
	return &Packager{
		s3:       s3,
		sts:      sts,
		region:   region,
		s3Suffix: s3Suffix,
		s3Tags:   s3Tags,
	}
}

type TemplateReader interface {
	fmt.Stringer
	Name() string
}

func (p *Packager) Package(ctx context.Context, prefix string, templateReader TemplateReader, writer io.Writer) (*string, error) {
	c, err := cfnyaml.Parse([]byte(templateReader.String()))
	if err != nil {
		return nil, err
	}

	nodes, err := c.PackageableNodes()
	artifacts := map[string]string{}
	for _, n := range nodes {
		path := n.Value
		realPath := filepath.Join(filepath.Dir(templateReader.Name()), path)
		if err != nil {
			return nil, errors.Wrap(err, "determining artifact path relative to template")
		}

		artifacts[path], err = zipper.Zip(realPath)
		if err != nil {
			return nil, errors.Wrapf(err, "zipping `%s`", path)
		}
	}

	uploads := map[string]*UploadedObject{}
	errch := make(chan error)
	for artifactPath, zipPath := range artifacts {
		go func(artifactPath, zipPath string) {
			hash, _ := md5path(zipPath)
			basename := strings.TrimSuffix(filepath.Base(artifactPath), ".zip")
			key := strings.TrimPrefix(fmt.Sprintf("%s/%s.zip/%s", prefix, basename, hash), "/")
			up, err := p.Upload(ctx, key, zipPath)
			uploads[artifactPath] = up
			errch <- errors.Wrap(err, "uploading zip to s3")
			if up != nil {
				if up.AlreadyExists {
					fmt.Fprintf(writer, "%s already exists at s3://%s/%s (v = %s)\n", artifactPath, up.Bucket, up.Key, up.VersionId)
				} else {
					fmt.Fprintf(writer, "Uploaded %s to s3://%s/%s (v = %s)\n", artifactPath, up.Bucket, up.Key, up.VersionId)
				}
			}
		}(artifactPath, zipPath)
	}

	for range artifacts {
		err = <-errch
		if err != nil {
			return nil, err
		}
	}

	for _, n := range nodes {
		path := n.Value
		uploaded := uploads[path]
		n.Replace(uploaded.Bucket, uploaded.Key, uploaded.VersionId)
	}

	templateBody := c.String()
	return &templateBody, nil
}

func md5path(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", errors.WithStack(err)
	}

	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", errors.WithStack(err)
	}

	sum := h.Sum(nil)
	return hex.EncodeToString(sum), nil
}
