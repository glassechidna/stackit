package stackit

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/pkg/errors"
	"time"
)

func (s *Stackit) Transform(template string, paramMap map[string]string) (*string, error) {
	params := []*cloudformation.Parameter{}
	for name, value := range paramMap {
		params = append(params, &cloudformation.Parameter{
			ParameterKey:   aws.String(name),
			ParameterValue: aws.String(value),
		})
	}

	stackName := fmt.Sprintf("stackit-temp-%d", time.Now().Unix())

	createResp, err := s.api.CreateChangeSet(&cloudformation.CreateChangeSetInput{
		ChangeSetName: aws.String(fmt.Sprintf("csid-%d", time.Now().Unix())),
		StackName:     &stackName,
		Capabilities:  aws.StringSlice([]string{"CAPABILITY_IAM", "CAPABILITY_NAMED_IAM"}),
		TemplateBody:  &template,
		ChangeSetType: aws.String(cloudformation.ChangeSetTypeCreate),
		Parameters:    params,
	})
	if err != nil {
		return nil, errors.Wrap(err, "creating change set")
	}

	_, err = s.waitForChangeset(createResp.Id)
	if err != nil {
		return nil, errors.Wrap(err, "waiting for change set")
	}

	getResp, err := s.api.GetTemplate(&cloudformation.GetTemplateInput{
		ChangeSetName: createResp.Id,
		TemplateStage: aws.String(cloudformation.TemplateStageProcessed),
	})
	if err != nil {
		return nil, errors.Wrap(err, "getting template body")
	}

	_, err = s.api.DeleteStack(&cloudformation.DeleteStackInput{StackName: &stackName})
	if err != nil {
		return nil, errors.Wrap(err, "deleting temporary stack")
	}

	return getResp.TemplateBody, err
}

func (s *Stackit) waitForChangeset(id *string) (*cloudformation.DescribeChangeSetOutput, error) {
	status := "CREATE_PENDING"
	terminal := []string{"CREATE_COMPLETE", "DELETE_COMPLETE", "FAILED"}

	var resp *cloudformation.DescribeChangeSetOutput
	var err error

	for !stringInSlice(terminal, status) {
		resp, err = s.api.DescribeChangeSet(&cloudformation.DescribeChangeSetInput{
			ChangeSetName: id,
		})
		if err != nil {
			return resp, errors.Wrap(err, "describing change set")
		}

		status = *resp.Status
		reason := ""
		if resp.StatusReason != nil {
			reason = *resp.StatusReason
		}

		if status == "FAILED" && reason != "The submitted information didn't contain changes. Submit different information to create a change set." {
			return resp, errors.New(reason)
		}

		time.Sleep(2 * time.Second)
	}

	return resp, nil
}
