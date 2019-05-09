package stackit

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Packager struct {
	s3     s3iface.S3API
	cfn    cloudformationiface.CloudFormationAPI
	sts    stsiface.STSAPI
	region string

	cachedBucketName string
}

func NewPackager(s3 s3iface.S3API, sts stsiface.STSAPI, region string) *Packager {
	return &Packager{
		s3:     s3,
		sts:    sts,
		region: region,
	}
}

type UploadedObject struct {
	Bucket    string
	Key       string
	VersionId string
}

func (p *Packager) Package(stackName, templatePath string, tags, parameters map[string]string) (*StackitUpInput, error) {
	absPath, err := filepath.Abs(templatePath)
	if err != nil {
		return nil, errors.Wrapf(err, "determining absolute path of '%s'", templatePath)
	}

	bucket, err := p.s3BucketName()
	if err != nil {
		return nil, err
	}

	prefix := stackName
	cliArgs := []string{"aws", "cloudformation", "package", "--template-file", absPath, "--s3-bucket", bucket, "--s3-prefix", prefix}
	cmd := exec.Command(cliArgs[0], cliArgs[1:]...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.Wrapf(err, "error running '%s': %s", strings.Join(cliArgs, " "), string(output))
	}

	templateBody := extractTemplateFromCliStdout(output)

	cfnParams := []*cloudformation.Parameter{}
	for name, value := range parameters {
		cfnParams = append(cfnParams, &cloudformation.Parameter{
			ParameterKey:   aws.String(name),
			ParameterValue: aws.String(value),
		})
	}

	upInput := StackitUpInput{
		StackName:       stackName,
		TemplateBody:    templateBody,
		PopulateMissing: true,
		Tags:            tags,
		Parameters:      cfnParams,
	}

	return &upInput, nil
}

func extractTemplateFromCliStdout(input []byte) string {
	lines := strings.Split(string(input), "\n")
	idx := 0
	for idx < len(lines) {
		if !strings.HasPrefix(lines[idx], "Uploading") {
			break
		}
		idx++
	}

	return strings.Join(lines[idx:], "\n")
}

func (p *Packager) Upload(prefix, path string) (*UploadedObject, error) {
	uploader := s3manager.NewUploaderWithClient(p.s3)

	bucket, err := p.s3BucketName()
	if err != nil {
		return nil, errors.Wrap(err, "getting bucket for object upload")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, errors.Wrapf(err, "determining absolute path of '%s'", path)
	}
	key := fmt.Sprintf("%s/%s", strings.Trim(prefix, "/"), filepath.Base(absPath))

	fi, err := os.Stat(absPath)
	if err != nil {
		return nil, errors.Wrapf(err, "determining if '%s' is a directory", path)
	}

	if fi.IsDir() {
		absPath = p.makeZip(absPath)
	}

	file, err := os.Open(absPath)
	if err != nil {
		return nil, errors.Wrapf(err, "opening file '%s'", absPath)
	}

	resp, err := uploader.Upload(&s3manager.UploadInput{
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

func (p *Packager) makeZip(path string) string {
	panic(errors.New("not yet implemented"))
}
