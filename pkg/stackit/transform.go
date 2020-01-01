package stackit

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/glassechidna/stackit/pkg/stackit/changeset"
	"github.com/pkg/errors"
	"time"
)

func (s *Stackit) Transform(ctx context.Context, template string, paramMap map[string]string) (*string, error) {
	params := []*cloudformation.Parameter{}
	for name, value := range paramMap {
		params = append(params, &cloudformation.Parameter{
			ParameterKey:   aws.String(name),
			ParameterValue: aws.String(value),
		})
	}

	stackName := fmt.Sprintf("stackit-temp-%d", time.Now().Unix())

	createResp, err := s.api.CreateChangeSetWithContext(ctx, &cloudformation.CreateChangeSetInput{
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

	_, err = changeset.Wait(ctx, s.api, *createResp.Id)
	if err != nil {
		return nil, errors.Wrap(err, "waiting for change set")
	}

	getResp, err := s.api.GetTemplateWithContext(ctx, &cloudformation.GetTemplateInput{
		ChangeSetName: createResp.Id,
		TemplateStage: aws.String(cloudformation.TemplateStageProcessed),
	})
	if err != nil {
		return nil, errors.Wrap(err, "getting template body")
	}

	_, err = s.api.DeleteStackWithContext(ctx, &cloudformation.DeleteStackInput{StackName: &stackName})
	if err != nil {
		return nil, errors.Wrap(err, "deleting temporary stack")
	}

	return getResp.TemplateBody, err
}
