package stackit

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/glassechidna/awsctx/service/cloudformationctx"
	"github.com/glassechidna/awsctx/service/stsctx"
	"github.com/pkg/errors"
	"io"
	"log"
)

type Stackit struct {
	api    cloudformationctx.CloudFormation
	stsApi stsctx.STS
}

func NewStackit(api cloudformationctx.CloudFormation, stsApi stsctx.STS) *Stackit {
	return &Stackit{api: api, stsApi: stsApi}
}

func (s *Stackit) Describe(ctx context.Context, stackName string) (*cloudformation.Stack, error) {
	resp, err := s.api.DescribeStacksWithContext(ctx, &cloudformation.DescribeStacksInput{StackName: &stackName})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			code := awsErr.Code()
			if code == "ThrottlingException" {
				return s.Describe(ctx, stackName)
			} else if code == "ValidationError" {
				return nil, nil
			}
		}
		return nil, errors.Wrap(err, "determining stack status")
	}

	stack := resp.Stacks[0]
	return stack, nil
}

func (s *Stackit) PrintOutputs(ctx context.Context, stackName string, writer io.Writer) {
	stack, err := s.Describe(ctx, stackName)

	if err != nil {
		log.Fatal(err.Error())
	}

	outputMap := make(map[string]string)

	for _, output := range stack.Outputs {
		outputMap[*output.OutputKey] = *output.OutputValue
	}

	bytes, err := json.MarshalIndent(outputMap, "", "  ")
	fmt.Fprintln(writer, string(bytes))
}

func (s *Stackit) IsSuccessfulState(ctx context.Context, stackName string) (bool, error) {
	stack, err := s.Describe(ctx, stackName)
	if err != nil {
		return false, errors.Wrap(err, "determining stack status")
	}

	status := *stack.StackStatus
	return status == "CREATE_COMPLETE" || status == "UPDATE_COMPLETE", nil
}
