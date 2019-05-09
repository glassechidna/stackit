package stackit

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/pkg/errors"
	"io"
	"log"
)

type Stackit struct {
	api    cloudformationiface.CloudFormationAPI
	stsApi stsiface.STSAPI
}

func NewStackit(api cloudformationiface.CloudFormationAPI, stsApi stsiface.STSAPI) *Stackit {
	return &Stackit{api: api, stsApi: stsApi}
}

func (s *Stackit) Describe(stackName string) (*cloudformation.Stack, error) {
	resp, err := s.api.DescribeStacks(&cloudformation.DescribeStacksInput{StackName: &stackName})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			code := awsErr.Code()
			if code == "ThrottlingException" {
				return s.Describe(stackName)
			} else if code == "ValidationError" {
				return nil, nil
			}
		}
		return nil, errors.Wrap(err, "determining stack status")
	}

	stack := resp.Stacks[0]
	return stack, nil
}

func (s *Stackit) PrintOutputs(stackName string, writer io.Writer) {
	stack, err := s.Describe(stackName)

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

func (s *Stackit) IsSuccessfulState(stackName string) (bool, error) {
	stack, err := s.Describe(stackName)
	if err != nil {
		return false, errors.Wrap(err, "determining stack status")
	}

	status := *stack.StackStatus
	return status == "CREATE_COMPLETE" || status == "UPDATE_COMPLETE", nil
}
