package kms

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/hashicorp/aws-sdk-go-base/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	awspolicy "github.com/jen20/awspolicyequivalence"
	"github.com/terraform-providers/terraform-provider-aws/internal/keyvaluetags"
	tfiam "github.com/terraform-providers/terraform-provider-aws/internal/service/iam"
	"github.com/terraform-providers/terraform-provider-aws/internal/tfresource"
)

const (
	// Maximum amount of time to wait for statusKeyState to return PendingDeletion
	keyStatePendingDeletionTimeout = 20 * time.Minute

	keyDeletedTimeout                = 20 * time.Minute
	keyDescriptionPropagationTimeout = 5 * time.Minute
	keyMaterialImportedTimeout       = 10 * time.Minute
	keyRotationUpdatedTimeout        = 10 * time.Minute
	keyStatePropagationTimeout       = 20 * time.Minute
	keyValidToPropagationTimeout     = 5 * time.Minute

	propagationTimeout = 2 * time.Minute
)

// waitIAMPropagation retries the specified function if the returned error indicates an IAM eventual consistency issue.
// If the retries time out the specified function is called one last time.
func waitIAMPropagation(f func() (interface{}, error)) (interface{}, error) {
	return tfresource.RetryWhenAwsErrCodeEquals(tfiam.PropagationTimeout, f, kms.ErrCodeMalformedPolicyDocumentException)
}

func waitKeyDeleted(conn *kms.KMS, id string) (*kms.KeyMetadata, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{kms.KeyStateDisabled, kms.KeyStateEnabled},
		Target:  []string{},
		Refresh: statusKeyState(conn, id),
		Timeout: keyDeletedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*kms.KeyMetadata); ok {
		return output, err
	}

	return nil, err
}

func waitKeyDescriptionPropagated(conn *kms.KMS, id string, description string) error {
	checkFunc := func() (bool, error) {
		output, err := findKeyByID(conn, id)

		if tfresource.NotFound(err) {
			return false, nil
		}

		if err != nil {
			return false, err
		}

		return aws.StringValue(output.Description) == description, nil
	}
	opts := tfresource.WaitOpts{
		ContinuousTargetOccurence: 5,
		MinTimeout:                2 * time.Second,
	}

	return tfresource.WaitUntil(keyDescriptionPropagationTimeout, checkFunc, opts)
}

func waitKeyMaterialImported(conn *kms.KMS, id string) (*kms.KeyMetadata, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{kms.KeyStatePendingImport},
		Target:  []string{kms.KeyStateDisabled, kms.KeyStateEnabled},
		Refresh: statusKeyState(conn, id),
		Timeout: keyMaterialImportedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*kms.KeyMetadata); ok {
		return output, err
	}

	return nil, err
}

func waitKeyPolicyPropagated(conn *kms.KMS, id, policy string) error {
	checkFunc := func() (bool, error) {
		output, err := findKeyPolicyByKeyIDAndPolicyName(conn, id, policyNameDefault)

		if tfresource.NotFound(err) {
			return false, nil
		}

		if err != nil {
			return false, err
		}

		equivalent, err := awspolicy.PoliciesAreEquivalent(aws.StringValue(output), policy)

		if err != nil {
			return false, err
		}

		return equivalent, nil
	}
	opts := tfresource.WaitOpts{
		ContinuousTargetOccurence: 5,
		MinTimeout:                1 * time.Second,
	}

	return tfresource.WaitUntil(propagationTimeout, checkFunc, opts)
}

func waitKeyRotationEnabledPropagated(conn *kms.KMS, id string, enabled bool) error {
	checkFunc := func() (bool, error) {
		output, err := findKeyRotationEnabledByKeyID(conn, id)

		if tfresource.NotFound(err) {
			return false, nil
		}

		if err != nil {
			return false, err
		}

		return aws.BoolValue(output) == enabled, nil
	}
	opts := tfresource.WaitOpts{
		ContinuousTargetOccurence: 5,
		MinTimeout:                1 * time.Second,
	}

	return tfresource.WaitUntil(propagationTimeout, checkFunc, opts)
}

func waitKeyStatePropagated(conn *kms.KMS, id string, enabled bool) error {
	checkFunc := func() (bool, error) {
		output, err := findKeyByID(conn, id)

		if tfresource.NotFound(err) {
			return false, nil
		}

		if err != nil {
			return false, err
		}

		return aws.BoolValue(output.Enabled) == enabled, nil
	}
	opts := tfresource.WaitOpts{
		ContinuousTargetOccurence: 15,
		MinTimeout:                2 * time.Second,
	}

	return tfresource.WaitUntil(keyStatePropagationTimeout, checkFunc, opts)
}

func waitKeyValidToPropagated(conn *kms.KMS, id string, validTo string) error {
	checkFunc := func() (bool, error) {
		output, err := findKeyByID(conn, id)

		if tfresource.NotFound(err) {
			return false, nil
		}

		if err != nil {
			return false, err
		}

		if output.ValidTo != nil {
			return aws.TimeValue(output.ValidTo).Format(time.RFC3339) == validTo, nil
		}

		return validTo == "", nil
	}
	opts := tfresource.WaitOpts{
		ContinuousTargetOccurence: 5,
		MinTimeout:                2 * time.Second,
	}

	return tfresource.WaitUntil(keyValidToPropagationTimeout, checkFunc, opts)
}

func waitTagsPropagated(conn *kms.KMS, id string, tags keyvaluetags.KeyValueTags) error {
	checkFunc := func() (bool, error) {
		output, err := keyvaluetags.KmsListTags(conn, id)

		if tfawserr.ErrCodeEquals(err, kms.ErrCodeNotFoundException) {
			return false, nil
		}

		if err != nil {
			return false, err
		}

		return output.Equal(tags), nil
	}
	opts := tfresource.WaitOpts{
		ContinuousTargetOccurence: 5,
		MinTimeout:                1 * time.Second,
	}

	return tfresource.WaitUntil(propagationTimeout, checkFunc, opts)
}
