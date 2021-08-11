package sagemaker

import (
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sagemaker"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-aws/internal/tfresource"
)

const (
	notebookInstanceInServiceTimeout  = 60 * time.Minute
	notebookInstanceStoppedTimeout    = 10 * time.Minute
	notebookInstanceDeletedTimeout    = 10 * time.Minute
	modelPackageGroupCompletedTimeout = 10 * time.Minute
	modelPackageGroupDeletedTimeout   = 10 * time.Minute
	imageCreatedTimeout               = 10 * time.Minute
	imageDeletedTimeout               = 10 * time.Minute
	imageVersionCreatedTimeout        = 10 * time.Minute
	imageVersionDeletedTimeout        = 10 * time.Minute
	domainInServiceTimeout            = 10 * time.Minute
	domainDeletedTimeout              = 10 * time.Minute
	featureGroupCreatedTimeout        = 10 * time.Minute
	featureGroupDeletedTimeout        = 10 * time.Minute
	userProfileInServiceTimeout       = 10 * time.Minute
	userProfileDeletedTimeout         = 10 * time.Minute
	appInServiceTimeout               = 10 * time.Minute
	appDeletedTimeout                 = 10 * time.Minute
)

// waitNotebookInstanceInService waits for a NotebookInstance to return InService
func waitNotebookInstanceInService(conn *sagemaker.SageMaker, notebookName string) (*sagemaker.DescribeNotebookInstanceOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			sagemakerNotebookInstanceStatusNotFound,
			sagemaker.NotebookInstanceStatusUpdating,
			sagemaker.NotebookInstanceStatusPending,
			sagemaker.NotebookInstanceStatusStopped,
		},
		Target:  []string{sagemaker.NotebookInstanceStatusInService},
		Refresh: statusNotebookInstance(conn, notebookName),
		Timeout: notebookInstanceInServiceTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*sagemaker.DescribeNotebookInstanceOutput); ok {
		return output, err
	}

	return nil, err
}

// waitNotebookInstanceStopped waits for a NotebookInstance to return Stopped
func waitNotebookInstanceStopped(conn *sagemaker.SageMaker, notebookName string) (*sagemaker.DescribeNotebookInstanceOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			sagemaker.NotebookInstanceStatusUpdating,
			sagemaker.NotebookInstanceStatusStopping,
		},
		Target:  []string{sagemaker.NotebookInstanceStatusStopped},
		Refresh: statusNotebookInstance(conn, notebookName),
		Timeout: notebookInstanceStoppedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*sagemaker.DescribeNotebookInstanceOutput); ok {
		return output, err
	}

	return nil, err
}

// waitNotebookInstanceDeleted waits for a NotebookInstance to return Deleted
func waitNotebookInstanceDeleted(conn *sagemaker.SageMaker, notebookName string) (*sagemaker.DescribeNotebookInstanceOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			sagemaker.NotebookInstanceStatusDeleting,
		},
		Target:  []string{},
		Refresh: statusNotebookInstance(conn, notebookName),
		Timeout: notebookInstanceDeletedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*sagemaker.DescribeNotebookInstanceOutput); ok {
		return output, err
	}

	return nil, err
}

// waitModelPackageGroupCompleted waits for a ModelPackageGroup to return Created
func waitModelPackageGroupCompleted(conn *sagemaker.SageMaker, name string) (*sagemaker.DescribeModelPackageGroupOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			sagemaker.ModelPackageGroupStatusPending,
			sagemaker.ModelPackageGroupStatusInProgress,
		},
		Target:  []string{sagemaker.ModelPackageGroupStatusCompleted},
		Refresh: statusModelPackageGroup(conn, name),
		Timeout: modelPackageGroupCompletedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*sagemaker.DescribeModelPackageGroupOutput); ok {
		return output, err
	}

	return nil, err
}

// waitModelPackageGroupDeleted waits for a ModelPackageGroup to return Created
func waitModelPackageGroupDeleted(conn *sagemaker.SageMaker, name string) (*sagemaker.DescribeModelPackageGroupOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			sagemaker.ModelPackageGroupStatusDeleting,
		},
		Target:  []string{},
		Refresh: statusModelPackageGroup(conn, name),
		Timeout: modelPackageGroupDeletedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*sagemaker.DescribeModelPackageGroupOutput); ok {
		return output, err
	}

	return nil, err
}

// waitImageCreated waits for a Image to return Created
func waitImageCreated(conn *sagemaker.SageMaker, name string) (*sagemaker.DescribeImageOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			sagemaker.ImageStatusCreating,
			sagemaker.ImageStatusUpdating,
		},
		Target:  []string{sagemaker.ImageStatusCreated},
		Refresh: statusImage(conn, name),
		Timeout: imageCreatedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*sagemaker.DescribeImageOutput); ok {
		return output, err
	}

	return nil, err
}

// waitImageDeleted waits for a Image to return Deleted
func waitImageDeleted(conn *sagemaker.SageMaker, name string) (*sagemaker.DescribeImageOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{sagemaker.ImageStatusDeleting},
		Target:  []string{},
		Refresh: statusImage(conn, name),
		Timeout: imageDeletedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*sagemaker.DescribeImageOutput); ok {
		return output, err
	}

	return nil, err
}

// waitImageVersionCreated waits for a ImageVersion to return Created
func waitImageVersionCreated(conn *sagemaker.SageMaker, name string) (*sagemaker.DescribeImageVersionOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			sagemaker.ImageVersionStatusCreating,
		},
		Target:  []string{sagemaker.ImageVersionStatusCreated},
		Refresh: statusImageVersion(conn, name),
		Timeout: imageVersionCreatedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*sagemaker.DescribeImageVersionOutput); ok {
		return output, err
	}

	return nil, err
}

// waitImageVersionDeleted waits for a ImageVersion to return Deleted
func waitImageVersionDeleted(conn *sagemaker.SageMaker, name string) (*sagemaker.DescribeImageVersionOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{sagemaker.ImageVersionStatusDeleting},
		Target:  []string{},
		Refresh: statusImageVersion(conn, name),
		Timeout: imageVersionDeletedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*sagemaker.DescribeImageVersionOutput); ok {
		return output, err
	}

	return nil, err
}

// waitDomainInService waits for a Domain to return InService
func waitDomainInService(conn *sagemaker.SageMaker, domainID string) (*sagemaker.DescribeDomainOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			sagemakerDomainStatusNotFound,
			sagemaker.DomainStatusPending,
			sagemaker.DomainStatusUpdating,
		},
		Target:  []string{sagemaker.DomainStatusInService},
		Refresh: statusDomain(conn, domainID),
		Timeout: domainInServiceTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*sagemaker.DescribeDomainOutput); ok {
		return output, err
	}

	return nil, err
}

// waitDomainDeleted waits for a Domain to return Deleted
func waitDomainDeleted(conn *sagemaker.SageMaker, domainID string) (*sagemaker.DescribeDomainOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			sagemaker.DomainStatusDeleting,
		},
		Target:  []string{},
		Refresh: statusDomain(conn, domainID),
		Timeout: domainDeletedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*sagemaker.DescribeDomainOutput); ok {
		return output, err
	}

	return nil, err
}

// waitFeatureGroupCreated waits for a Feature Group to return Created
func waitFeatureGroupCreated(conn *sagemaker.SageMaker, name string) (*sagemaker.DescribeFeatureGroupOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{sagemaker.FeatureGroupStatusCreating},
		Target:  []string{sagemaker.FeatureGroupStatusCreated},
		Refresh: statusFeatureGroup(conn, name),
		Timeout: featureGroupCreatedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*sagemaker.DescribeFeatureGroupOutput); ok {
		if status, reason := aws.StringValue(output.FeatureGroupStatus), aws.StringValue(output.FailureReason); status == sagemaker.FeatureGroupStatusCreateFailed && reason != "" {
			tfresource.SetLastError(err, errors.New(reason))
		}

		return output, err
	}

	return nil, err
}

// waitFeatureGroupDeleted waits for a Feature Group to return Deleted
func waitFeatureGroupDeleted(conn *sagemaker.SageMaker, name string) (*sagemaker.DescribeFeatureGroupOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{sagemaker.FeatureGroupStatusDeleting},
		Target:  []string{},
		Refresh: statusFeatureGroup(conn, name),
		Timeout: featureGroupDeletedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*sagemaker.DescribeFeatureGroupOutput); ok {
		if status, reason := aws.StringValue(output.FeatureGroupStatus), aws.StringValue(output.FailureReason); status == sagemaker.FeatureGroupStatusDeleteFailed && reason != "" {
			tfresource.SetLastError(err, errors.New(reason))
		}

		return output, err
	}

	return nil, err
}

// waitUserProfileInService waits for a UserProfile to return InService
func waitUserProfileInService(conn *sagemaker.SageMaker, domainID, userProfileName string) (*sagemaker.DescribeUserProfileOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			sagemakerUserProfileStatusNotFound,
			sagemaker.UserProfileStatusPending,
			sagemaker.UserProfileStatusUpdating,
		},
		Target:  []string{sagemaker.UserProfileStatusInService},
		Refresh: statusUserProfile(conn, domainID, userProfileName),
		Timeout: userProfileInServiceTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*sagemaker.DescribeUserProfileOutput); ok {
		return output, err
	}

	return nil, err
}

// waitUserProfileDeleted waits for a UserProfile to return Deleted
func waitUserProfileDeleted(conn *sagemaker.SageMaker, domainID, userProfileName string) (*sagemaker.DescribeUserProfileOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			sagemaker.UserProfileStatusDeleting,
		},
		Target:  []string{},
		Refresh: statusUserProfile(conn, domainID, userProfileName),
		Timeout: userProfileDeletedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*sagemaker.DescribeUserProfileOutput); ok {
		return output, err
	}

	return nil, err
}

// waitAppInService waits for a App to return InService
func waitAppInService(conn *sagemaker.SageMaker, domainID, userProfileName, appType, appName string) (*sagemaker.DescribeAppOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			sagemakerAppStatusNotFound,
			sagemaker.AppStatusPending,
		},
		Target:  []string{sagemaker.AppStatusInService},
		Refresh: statusApp(conn, domainID, userProfileName, appType, appName),
		Timeout: appInServiceTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*sagemaker.DescribeAppOutput); ok {
		return output, err
	}

	return nil, err
}

// waitAppDeleted waits for a App to return Deleted
func waitAppDeleted(conn *sagemaker.SageMaker, domainID, userProfileName, appType, appName string) (*sagemaker.DescribeAppOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			sagemaker.AppStatusDeleting,
		},
		Target: []string{
			sagemaker.AppStatusDeleted,
		},
		Refresh: statusApp(conn, domainID, userProfileName, appType, appName),
		Timeout: appDeletedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*sagemaker.DescribeAppOutput); ok {
		return output, err
	}

	return nil, err
}
