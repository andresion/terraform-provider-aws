package sagemaker

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sagemaker"
	"github.com/hashicorp/aws-sdk-go-base/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-aws/internal/tfresource"
)

const (
	sagemakerNotebookInstanceStatusNotFound  = "NotFound"
	sagemakerImageStatusNotFound             = "NotFound"
	sagemakerImageStatusFailed               = "Failed"
	sagemakerImageVersionStatusNotFound      = "NotFound"
	sagemakerImageVersionStatusFailed        = "Failed"
	sagemakerDomainStatusNotFound            = "NotFound"
	sagemakerUserProfileStatusNotFound       = "NotFound"
	sagemakerModelPackageGroupStatusNotFound = "NotFound"
	sagemakerAppStatusNotFound               = "NotFound"
)

// statusNotebookInstance fetches the NotebookInstance and its Status
func statusNotebookInstance(conn *sagemaker.SageMaker, notebookName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		input := &sagemaker.DescribeNotebookInstanceInput{
			NotebookInstanceName: aws.String(notebookName),
		}

		output, err := conn.DescribeNotebookInstance(input)

		if tfawserr.ErrMessageContains(err, "ValidationException", "RecordNotFound") {
			return nil, sagemakerNotebookInstanceStatusNotFound, nil
		}

		if err != nil {
			return nil, sagemaker.NotebookInstanceStatusFailed, err
		}

		if output == nil {
			return nil, sagemakerNotebookInstanceStatusNotFound, nil
		}

		return output, aws.StringValue(output.NotebookInstanceStatus), nil
	}
}

// statusModelPackageGroup fetches the ModelPackageGroup and its Status
func statusModelPackageGroup(conn *sagemaker.SageMaker, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		input := &sagemaker.DescribeModelPackageGroupInput{
			ModelPackageGroupName: aws.String(name),
		}

		output, err := conn.DescribeModelPackageGroup(input)

		if tfawserr.ErrMessageContains(err, "ValidationException", "does not exist") {
			return nil, sagemakerModelPackageGroupStatusNotFound, nil
		}

		if err != nil {
			return nil, sagemaker.ModelPackageGroupStatusFailed, err
		}

		if output == nil {
			return nil, sagemakerModelPackageGroupStatusNotFound, nil
		}

		return output, aws.StringValue(output.ModelPackageGroupStatus), nil
	}
}

// statusImage fetches the Image and its Status
func statusImage(conn *sagemaker.SageMaker, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		input := &sagemaker.DescribeImageInput{
			ImageName: aws.String(name),
		}

		output, err := conn.DescribeImage(input)

		if tfawserr.ErrMessageContains(err, sagemaker.ErrCodeResourceNotFound, "No Image with the name") {
			return nil, sagemakerImageStatusNotFound, nil
		}

		if err != nil {
			return nil, sagemakerImageStatusFailed, err
		}

		if output == nil {
			return nil, sagemakerImageStatusNotFound, nil
		}

		if aws.StringValue(output.ImageStatus) == sagemaker.ImageStatusCreateFailed {
			return output, sagemaker.ImageStatusCreateFailed, fmt.Errorf("%s", aws.StringValue(output.FailureReason))
		}

		return output, aws.StringValue(output.ImageStatus), nil
	}
}

// statusImageVersion fetches the ImageVersion and its Status
func statusImageVersion(conn *sagemaker.SageMaker, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		input := &sagemaker.DescribeImageVersionInput{
			ImageName: aws.String(name),
		}

		output, err := conn.DescribeImageVersion(input)

		if tfawserr.ErrMessageContains(err, sagemaker.ErrCodeResourceNotFound, "No ImageVersion with the name") {
			return nil, sagemakerImageVersionStatusNotFound, nil
		}

		if err != nil {
			return nil, sagemakerImageVersionStatusFailed, err
		}

		if output == nil {
			return nil, sagemakerImageVersionStatusNotFound, nil
		}

		if aws.StringValue(output.ImageVersionStatus) == sagemaker.ImageVersionStatusCreateFailed {
			return output, sagemaker.ImageVersionStatusCreateFailed, fmt.Errorf("%s", aws.StringValue(output.FailureReason))
		}

		return output, aws.StringValue(output.ImageVersionStatus), nil
	}
}

// statusDomain fetches the Domain and its Status
func statusDomain(conn *sagemaker.SageMaker, domainID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		input := &sagemaker.DescribeDomainInput{
			DomainId: aws.String(domainID),
		}

		output, err := conn.DescribeDomain(input)

		if tfawserr.ErrMessageContains(err, "ValidationException", "RecordNotFound") {
			return nil, sagemaker.UserProfileStatusFailed, nil
		}

		if err != nil {
			return nil, sagemaker.DomainStatusFailed, err
		}

		if output == nil {
			return nil, sagemakerDomainStatusNotFound, nil
		}

		return output, aws.StringValue(output.Status), nil
	}
}

func statusFeatureGroup(conn *sagemaker.SageMaker, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := findFeatureGroupByName(conn, name)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.FeatureGroupStatus), nil
	}
}

// statusUserProfile fetches the UserProfile and its Status
func statusUserProfile(conn *sagemaker.SageMaker, domainID, userProfileName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		input := &sagemaker.DescribeUserProfileInput{
			DomainId:        aws.String(domainID),
			UserProfileName: aws.String(userProfileName),
		}

		output, err := conn.DescribeUserProfile(input)

		if tfawserr.ErrMessageContains(err, "ValidationException", "RecordNotFound") {
			return nil, sagemakerUserProfileStatusNotFound, nil
		}

		if err != nil {
			return nil, sagemaker.UserProfileStatusFailed, err
		}

		if output == nil {
			return nil, sagemakerUserProfileStatusNotFound, nil
		}

		return output, aws.StringValue(output.Status), nil
	}
}

// statusApp fetches the App and its Status
func statusApp(conn *sagemaker.SageMaker, domainID, userProfileName, appType, appName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		input := &sagemaker.DescribeAppInput{
			DomainId:        aws.String(domainID),
			UserProfileName: aws.String(userProfileName),
			AppType:         aws.String(appType),
			AppName:         aws.String(appName),
		}

		output, err := conn.DescribeApp(input)

		if tfawserr.ErrMessageContains(err, "ValidationException", "RecordNotFound") {
			return nil, sagemakerAppStatusNotFound, nil
		}

		if err != nil {
			return nil, sagemaker.AppStatusFailed, err
		}

		if output == nil {
			return nil, sagemakerAppStatusNotFound, nil
		}

		return output, aws.StringValue(output.Status), nil
	}
}
