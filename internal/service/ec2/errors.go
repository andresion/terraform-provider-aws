package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	multierror "github.com/hashicorp/go-multierror"
)

const (
	errCodeGatewayNotAttached           = "Gateway.NotAttached"
	errCodeInvalidAssociationIDNotFound = "InvalidAssociationID.NotFound"
	errCodeInvalidParameter             = "InvalidParameter"
	errCodeInvalidParameterException    = "InvalidParameterException"
	errCodeInvalidParameterValue        = "InvalidParameterValue"
)

const (
	errCodeInvalidCarrierGatewayIDNotFound = "InvalidCarrierGatewayID.NotFound"
)

const (
	errCodeInvalidNetworkInterfaceIDNotFound = "InvalidNetworkInterfaceID.NotFound"
)

const (
	errCodeInvalidPrefixListIDNotFound = "InvalidPrefixListID.NotFound"
)

const (
	errCodeInvalidRouteNotFound        = "InvalidRoute.NotFound"
	errCodeInvalidRouteTableIdNotFound = "InvalidRouteTableId.NotFound"
	errCodeInvalidRouteTableIDNotFound = "InvalidRouteTableID.NotFound"
)

const (
	errCodeInvalidTransitGatewayIDNotFound = "InvalidTransitGatewayID.NotFound"
)

const (
	errCodeClientVPNEndpointIdNotFound        = "InvalidClientVpnEndpointId.NotFound"
	errCodeClientVPNAuthorizationRuleNotFound = "InvalidClientVpnEndpointAuthorizationRuleNotFound"
	errCodeClientVPNAssociationIdNotFound     = "InvalidClientVpnAssociationId.NotFound"
	errCodeClientVPNRouteNotFound             = "InvalidClientVpnRouteNotFound"
)

const (
	errCodeInvalidInstanceIDNotFound = "InvalidInstanceID.NotFound"
)

const (
	invalidSecurityGroupIDNotFound = "InvalidSecurityGroupID.NotFound"
	invalidGroupNotFound           = "InvalidGroup.NotFound"
)

const (
	errCodeInvalidSpotInstanceRequestIDNotFound = "InvalidSpotInstanceRequestID.NotFound"
)

const (
	errCodeInvalidSubnetIdNotFound = "InvalidSubnetId.NotFound"
	errCodeInvalidSubnetIDNotFound = "InvalidSubnetID.NotFound"
)

const (
	errCodeInvalidVPCIDNotFound = "InvalidVpcID.NotFound"
)

const (
	errCodeInvalidVPCEndpointIdNotFound        = "InvalidVpcEndpointId.NotFound"
	errCodeInvalidVPCEndpointNotFound          = "InvalidVpcEndpoint.NotFound"
	errCodeInvalidVPCEndpointServiceIdNotFound = "InvalidVpcEndpointServiceId.NotFound"
)

const (
	errCodeInvalidVPCPeeringConnectionIDNotFound = "InvalidVpcPeeringConnectionID.NotFound"
)

const (
	invalidVPNGatewayAttachmentNotFound = "InvalidVpnGatewayAttachment.NotFound"
	invalidVPNGatewayIDNotFound         = "InvalidVpnGatewayID.NotFound"
)

const (
	errCodeInvalidPermissionDuplicate = "InvalidPermission.Duplicate"
	errCodeInvalidPermissionMalformed = "InvalidPermission.Malformed"
	errCodeInvalidPermissionNotFound  = "InvalidPermission.NotFound"
)

// See https://docs.aws.amazon.com/vm-import/latest/userguide/vmimport-image-import.html#check-import-task-status
const (
	eBSSnapshotImportActive     = "active"
	eBSSnapshotImportDeleting   = "deleting"
	eBSSnapshotImportDeleted    = "deleted"
	eBSSnapshotImportUpdating   = "updating"
	eBSSnapshotImportValidating = "validating"
	eBSSnapshotImportValidated  = "validated"
	eBSSnapshotImportConverting = "converting"
	eBSSnapshotImportCompleted  = "completed"
)

func unsuccessfulItemError(apiObject *ec2.UnsuccessfulItemError) error {
	if apiObject == nil {
		return nil
	}

	return awserr.New(aws.StringValue(apiObject.Code), aws.StringValue(apiObject.Message), nil)
}

func unsuccessfulItemsError(apiObjects []*ec2.UnsuccessfulItem) error {
	var errors *multierror.Error

	for _, apiObject := range apiObjects {
		if apiObject == nil {
			continue
		}

		err := unsuccessfulItemError(apiObject.Error)

		if err != nil {
			errors = multierror.Append(errors, fmt.Errorf("%s: %w", aws.StringValue(apiObject.ResourceId), err))
		}
	}

	return errors.ErrorOrNil()
}
