package eks

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/eks"
	multierror "github.com/hashicorp/go-multierror"
)

func addonIssueError(apiObject *eks.AddonIssue) error {
	if apiObject == nil {
		return nil
	}

	return awserr.New(aws.StringValue(apiObject.Code), aws.StringValue(apiObject.Message), nil)
}

func addonIssuesError(apiObjects []*eks.AddonIssue) error {
	var errors *multierror.Error

	for _, apiObject := range apiObjects {
		if apiObject == nil {
			continue
		}

		err := addonIssueError(apiObject)

		if err != nil {
			errors = multierror.Append(errors, fmt.Errorf("%s: %w", strings.Join(aws.StringValueSlice(apiObject.ResourceIds), ", "), err))
		}
	}

	return errors.ErrorOrNil()
}

func errorDetailError(apiObject *eks.ErrorDetail) error {
	if apiObject == nil {
		return nil
	}

	return awserr.New(aws.StringValue(apiObject.ErrorCode), aws.StringValue(apiObject.ErrorMessage), nil)
}

func errorDetailsError(apiObjects []*eks.ErrorDetail) error {
	var errors *multierror.Error

	for _, apiObject := range apiObjects {
		if apiObject == nil {
			continue
		}

		err := errorDetailError(apiObject)

		if err != nil {
			errors = multierror.Append(errors, fmt.Errorf("%s: %w", strings.Join(aws.StringValueSlice(apiObject.ResourceIds), ", "), err))
		}
	}

	return errors.ErrorOrNil()
}

func issueError(apiObject *eks.Issue) error {
	if apiObject == nil {
		return nil
	}

	return awserr.New(aws.StringValue(apiObject.Code), aws.StringValue(apiObject.Message), nil)
}

func issuesError(apiObjects []*eks.Issue) error {
	var errors *multierror.Error

	for _, apiObject := range apiObjects {
		if apiObject == nil {
			continue
		}

		err := issueError(apiObject)

		if err != nil {
			errors = multierror.Append(errors, fmt.Errorf("%s: %w", strings.Join(aws.StringValueSlice(apiObject.ResourceIds), ", "), err))
		}
	}

	return errors.ErrorOrNil()
}
