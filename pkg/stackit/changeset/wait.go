package changeset

import (
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/pkg/errors"
	"time"
)

type NoOpChangesetError struct{}

const errNoOp = "The submitted information didn't contain changes. Submit different information to create a change set."
const errNoOp2 = "No updates are to be performed."

func (e *NoOpChangesetError) Error() string {
	return errNoOp
}

func Wait(api cloudformationiface.CloudFormationAPI, id string) (*cloudformation.DescribeChangeSetOutput, error) {
	status := "CREATE_PENDING"
	terminal := []string{"CREATE_COMPLETE", "DELETE_COMPLETE", "FAILED"}

	var resp *cloudformation.DescribeChangeSetOutput
	var err error

	for !stringInSlice(terminal, status) {
		resp, err = api.DescribeChangeSet(&cloudformation.DescribeChangeSetInput{
			ChangeSetName: &id,
		})
		if err != nil {
			return resp, errors.Wrap(err, "describing change set")
		}

		status = *resp.Status
		reason := ""
		if resp.StatusReason != nil {
			reason = *resp.StatusReason
		}

		if status == "FAILED" {
			if reason == errNoOp || reason == errNoOp2 {
				return resp, &NoOpChangesetError{}
			} else {
				return nil, errors.New(reason)
			}
		}

		time.Sleep(2 * time.Second)
	}

	return resp, nil
}

func stringInSlice(slice []string, s string) bool {
	for _, ss := range slice {
		if s == ss {
			return true
		}
	}
	return false
}
