package ram

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ram"
	"github.com/hashicorp/aws-sdk-go-base/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	resourceShareInvitationStatusNotFound = "NotFound"
	resourceShareInvitationStatusUnknown  = "Unknown"

	resourceShareStatusNotFound = "NotFound"
	resourceShareStatusUnknown  = "Unknown"

	principalAssociationStatusNotFound = "NotFound"
)

// statusResourceShareInvitation fetches the ResourceShareInvitation and its Status
func statusResourceShareInvitation(conn *ram.RAM, arn string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		invitation, err := findResourceShareInvitationByARN(conn, arn)

		if err != nil {
			return nil, resourceShareInvitationStatusUnknown, err
		}

		if invitation == nil {
			return nil, resourceShareInvitationStatusNotFound, nil
		}

		return invitation, aws.StringValue(invitation.Status), nil
	}
}

// statusResourceShareOwnerSelf fetches the ResourceShare and its Status
func statusResourceShareOwnerSelf(conn *ram.RAM, arn string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		share, err := findResourceShareOwnerSelfByARN(conn, arn)

		if tfawserr.ErrCodeEquals(err, ram.ErrCodeUnknownResourceException) {
			return nil, resourceShareStatusNotFound, nil
		}

		if err != nil {
			return nil, resourceShareStatusUnknown, err
		}

		if share == nil {
			return nil, resourceShareStatusNotFound, nil
		}

		return share, aws.StringValue(share.Status), nil
	}
}

func statusResourceSharePrincipalAssociation(conn *ram.RAM, resourceShareArn, principal string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		association, err := findResourceSharePrincipalAssociationByShareARNPrincipal(conn, resourceShareArn, principal)

		if tfawserr.ErrCodeEquals(err, ram.ErrCodeUnknownResourceException) {
			return nil, principalAssociationStatusNotFound, err
		}

		if err != nil {
			return nil, ram.ResourceShareAssociationStatusFailed, err
		}

		if association == nil {
			return nil, ram.ResourceShareAssociationStatusDisassociated, nil
		}

		if aws.StringValue(association.Status) == ram.ResourceShareAssociationStatusFailed {
			extendedErr := fmt.Errorf("association status message: %s", aws.StringValue(association.StatusMessage))
			return association, aws.StringValue(association.Status), extendedErr
		}

		return association, aws.StringValue(association.Status), nil
	}
}
