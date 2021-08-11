package ec2

import (
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/aws-sdk-go-base/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tfiam "github.com/terraform-providers/terraform-provider-aws/internal/service/iam"
	"github.com/terraform-providers/terraform-provider-aws/internal/tfresource"
)

const (
	carrierGatewayStateNotFound = "NotFound"
	carrierGatewayStateUnknown  = "Unknown"
	snapshotImportNotFound      = "NotFound"
)

// statusCarrierGatewayState fetches the CarrierGateway and its State
func statusCarrierGatewayState(conn *ec2.EC2, carrierGatewayID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		carrierGateway, err := findCarrierGatewayByID(conn, carrierGatewayID)
		if tfawserr.ErrCodeEquals(err, errCodeInvalidCarrierGatewayIDNotFound) {
			return nil, carrierGatewayStateNotFound, nil
		}
		if err != nil {
			return nil, carrierGatewayStateUnknown, err
		}

		if carrierGateway == nil {
			return nil, carrierGatewayStateNotFound, nil
		}

		state := aws.StringValue(carrierGateway.State)

		if state == ec2.CarrierGatewayStateDeleted {
			return nil, carrierGatewayStateNotFound, nil
		}

		return carrierGateway, state, nil
	}
}

// statusLocalGatewayRouteTableVPCAssociationState fetches the LocalGatewayRouteTableVpcAssociation and its State
func statusLocalGatewayRouteTableVPCAssociationState(conn *ec2.EC2, localGatewayRouteTableVpcAssociationID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		input := &ec2.DescribeLocalGatewayRouteTableVpcAssociationsInput{
			LocalGatewayRouteTableVpcAssociationIds: aws.StringSlice([]string{localGatewayRouteTableVpcAssociationID}),
		}

		output, err := conn.DescribeLocalGatewayRouteTableVpcAssociations(input)

		if err != nil {
			return nil, "", err
		}

		var association *ec2.LocalGatewayRouteTableVpcAssociation

		for _, outputAssociation := range output.LocalGatewayRouteTableVpcAssociations {
			if outputAssociation == nil {
				continue
			}

			if aws.StringValue(outputAssociation.LocalGatewayRouteTableVpcAssociationId) == localGatewayRouteTableVpcAssociationID {
				association = outputAssociation
				break
			}
		}

		if association == nil {
			return association, ec2.RouteTableAssociationStateCodeDisassociated, nil
		}

		return association, aws.StringValue(association.State), nil
	}
}

const (
	clientVPNEndpointStatusNotFound = "NotFound"

	clientVPNEndpointStatusUnknown = "Unknown"
)

// statusClientVPNEndpoint fetches the Client VPN endpoint and its Status
func statusClientVPNEndpoint(conn *ec2.EC2, endpointID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		result, err := conn.DescribeClientVpnEndpoints(&ec2.DescribeClientVpnEndpointsInput{
			ClientVpnEndpointIds: aws.StringSlice([]string{endpointID}),
		})
		if tfawserr.ErrCodeEquals(err, errCodeClientVPNEndpointIdNotFound) {
			return nil, clientVPNEndpointStatusNotFound, nil
		}
		if err != nil {
			return nil, clientVPNEndpointStatusUnknown, err
		}

		if result == nil || len(result.ClientVpnEndpoints) == 0 || result.ClientVpnEndpoints[0] == nil {
			return nil, clientVPNEndpointStatusNotFound, nil
		}

		endpoint := result.ClientVpnEndpoints[0]
		if endpoint.Status == nil || endpoint.Status.Code == nil {
			return endpoint, clientVPNEndpointStatusUnknown, nil
		}

		return endpoint, aws.StringValue(endpoint.Status.Code), nil
	}
}

const (
	clientVPNAuthorizationRuleStatusNotFound = "NotFound"

	clientVPNAuthorizationRuleStatusUnknown = "Unknown"
)

// statusClientVPNAuthorizationRule fetches the Client VPN authorization rule and its Status
func statusClientVPNAuthorizationRule(conn *ec2.EC2, authorizationRuleID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		result, err := findClientVPNAuthorizationRuleByID(conn, authorizationRuleID)
		if tfawserr.ErrCodeEquals(err, errCodeClientVPNAuthorizationRuleNotFound) {
			return nil, clientVPNAuthorizationRuleStatusNotFound, nil
		}
		if err != nil {
			return nil, clientVPNAuthorizationRuleStatusUnknown, err
		}

		if result == nil || len(result.AuthorizationRules) == 0 || result.AuthorizationRules[0] == nil {
			return nil, clientVPNAuthorizationRuleStatusNotFound, nil
		}

		if len(result.AuthorizationRules) > 1 {
			return nil, clientVPNAuthorizationRuleStatusUnknown, fmt.Errorf("internal error: found %d results for Client VPN authorization rule (%s) status, need 1", len(result.AuthorizationRules), authorizationRuleID)
		}

		rule := result.AuthorizationRules[0]
		if rule.Status == nil || rule.Status.Code == nil {
			return rule, clientVPNAuthorizationRuleStatusUnknown, nil
		}

		return rule, aws.StringValue(rule.Status.Code), nil
	}
}

const (
	clientVPNNetworkAssociationStatusNotFound = "NotFound"

	clientVPNNetworkAssociationStatusUnknown = "Unknown"
)

func statusClientVPNNetworkAssociation(conn *ec2.EC2, cvnaID string, cvepID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		result, err := conn.DescribeClientVpnTargetNetworks(&ec2.DescribeClientVpnTargetNetworksInput{
			ClientVpnEndpointId: aws.String(cvepID),
			AssociationIds:      []*string{aws.String(cvnaID)},
		})

		if tfawserr.ErrCodeEquals(err, errCodeClientVPNAssociationIdNotFound) || tfawserr.ErrCodeEquals(err, errCodeClientVPNEndpointIdNotFound) {
			return nil, clientVPNNetworkAssociationStatusNotFound, nil
		}
		if err != nil {
			return nil, clientVPNNetworkAssociationStatusUnknown, err
		}

		if result == nil || len(result.ClientVpnTargetNetworks) == 0 || result.ClientVpnTargetNetworks[0] == nil {
			return nil, clientVPNNetworkAssociationStatusNotFound, nil
		}

		network := result.ClientVpnTargetNetworks[0]
		if network.Status == nil || network.Status.Code == nil {
			return network, clientVPNNetworkAssociationStatusUnknown, nil
		}

		return network, aws.StringValue(network.Status.Code), nil
	}
}

const (
	clientVPNRouteStatusNotFound = "NotFound"

	clientVPNRouteStatusUnknown = "Unknown"
)

// statusClientVPNRoute fetches the Client VPN route and its Status
func statusClientVPNRoute(conn *ec2.EC2, routeID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		result, err := findClientVPNRouteByID(conn, routeID)
		if tfawserr.ErrCodeEquals(err, errCodeClientVPNRouteNotFound) {
			return nil, clientVPNRouteStatusNotFound, nil
		}
		if err != nil {
			return nil, clientVPNRouteStatusUnknown, err
		}

		if result == nil || len(result.Routes) == 0 || result.Routes[0] == nil {
			return nil, clientVPNRouteStatusNotFound, nil
		}

		if len(result.Routes) > 1 {
			return nil, clientVPNRouteStatusUnknown, fmt.Errorf("internal error: found %d results for Client VPN route (%s) status, need 1", len(result.Routes), routeID)
		}

		rule := result.Routes[0]
		if rule.Status == nil || rule.Status.Code == nil {
			return rule, clientVPNRouteStatusUnknown, nil
		}

		return rule, aws.StringValue(rule.Status.Code), nil
	}
}

// statusInstanceIAMInstanceProfile fetches the Instance and its IamInstanceProfile
//
// The EC2 API accepts a name and always returns an ARN, so it is converted
// back to the name to prevent unexpected differences.
func statusInstanceIAMInstanceProfile(conn *ec2.EC2, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := findInstanceByID(conn, id)

		if tfawserr.ErrCodeEquals(err, errCodeInvalidInstanceIDNotFound) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		if instance == nil {
			return nil, "", nil
		}

		if instance.IamInstanceProfile == nil || instance.IamInstanceProfile.Arn == nil {
			return instance, "", nil
		}

		name, err := tfiam.InstanceProfileARNToName(aws.StringValue(instance.IamInstanceProfile.Arn))

		if err != nil {
			return instance, "", err
		}

		return instance, name, nil
	}
}

const (
	routeStatusReady = "ready"
)

func statusRoute(conn *ec2.EC2, routeFinder routeFinder, routeTableID, destination string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := routeFinder(conn, routeTableID, destination)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, routeStatusReady, nil
	}
}

const (
	routeTableStatusReady = "ready"
)

func statusRouteTable(conn *ec2.EC2, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := findRouteTableByID(conn, id)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, routeTableStatusReady, nil
	}
}

func statusRouteTableAssociationState(conn *ec2.EC2, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := findRouteTableAssociationByID(conn, id)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output.AssociationState, aws.StringValue(output.AssociationState.State), nil
	}
}

const (
	securityGroupStatusCreated = "Created"

	securityGroupStatusNotFound = "NotFound"

	securityGroupStatusUnknown = "Unknown"
)

// statusSecurityGroup fetches the security group and its status
func statusSecurityGroup(conn *ec2.EC2, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		group, err := findSecurityGroupByID(conn, id)
		if tfresource.NotFound(err) {
			return nil, securityGroupStatusNotFound, nil
		}
		if err != nil {
			return nil, securityGroupStatusUnknown, err
		}

		return group, securityGroupStatusCreated, nil
	}
}

// statusSubnetMapCustomerOwnedIPOnLaunch fetches the Subnet and its MapCustomerOwnedIpOnLaunch
func statusSubnetMapCustomerOwnedIPOnLaunch(conn *ec2.EC2, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		subnet, err := findSubnetByID(conn, id)

		if tfawserr.ErrCodeEquals(err, errCodeInvalidSubnetIDNotFound) {
			return nil, "false", nil
		}

		if err != nil {
			return nil, "false", err
		}

		if subnet == nil {
			return nil, "false", nil
		}

		return subnet, strconv.FormatBool(aws.BoolValue(subnet.MapCustomerOwnedIpOnLaunch)), nil
	}
}

// statusSubnetMapPublicIPOnLaunch fetches the Subnet and its MapPublicIpOnLaunch
func statusSubnetMapPublicIPOnLaunch(conn *ec2.EC2, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		subnet, err := findSubnetByID(conn, id)

		if tfawserr.ErrCodeEquals(err, errCodeInvalidSubnetIDNotFound) {
			return nil, "false", nil
		}

		if err != nil {
			return nil, "false", err
		}

		if subnet == nil {
			return nil, "false", nil
		}

		return subnet, strconv.FormatBool(aws.BoolValue(subnet.MapPublicIpOnLaunch)), nil
	}
}

func statusTransitGatewayPrefixListReferenceState(conn *ec2.EC2, transitGatewayRouteTableID string, prefixListID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		transitGatewayPrefixListReference, err := findTransitGatewayPrefixListReference(conn, transitGatewayRouteTableID, prefixListID)

		if err != nil {
			return nil, "", err
		}

		if transitGatewayPrefixListReference == nil {
			return nil, "", nil
		}

		return transitGatewayPrefixListReference, aws.StringValue(transitGatewayPrefixListReference.State), nil
	}
}

func statusTransitGatewayRouteTablePropagationState(conn *ec2.EC2, transitGatewayRouteTableID string, transitGatewayAttachmentID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		transitGatewayRouteTablePropagation, err := findTransitGatewayRouteTablePropagation(conn, transitGatewayRouteTableID, transitGatewayAttachmentID)

		if err != nil {
			return nil, "", err
		}

		if transitGatewayRouteTablePropagation == nil {
			return nil, "", nil
		}

		return transitGatewayRouteTablePropagation, aws.StringValue(transitGatewayRouteTablePropagation.State), nil
	}
}

// statusVPCAttribute fetches the Vpc and its attribute value
func statusVPCAttribute(conn *ec2.EC2, id string, attribute string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		attributeValue, err := findVPCAttribute(conn, id, attribute)

		if tfawserr.ErrCodeEquals(err, errCodeInvalidVPCIDNotFound) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		if attributeValue == nil {
			return nil, "", nil
		}

		return attributeValue, strconv.FormatBool(aws.BoolValue(attributeValue)), nil
	}
}

const (
	vpcPeeringConnectionStatusNotFound = "NotFound"
	vpcPeeringConnectionStatusUnknown  = "Unknown"
)

// statusVPCPeeringConnection fetches the VPC peering connection and its status
func statusVPCPeeringConnection(conn *ec2.EC2, vpcPeeringConnectionID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		vpcPeeringConnection, err := findVPCPeeringConnectionByID(conn, vpcPeeringConnectionID)
		if tfawserr.ErrCodeEquals(err, errCodeInvalidVPCPeeringConnectionIDNotFound) {
			return nil, vpcPeeringConnectionStatusNotFound, nil
		}
		if err != nil {
			return nil, vpcPeeringConnectionStatusUnknown, err
		}

		// Sometimes AWS just has consistency issues and doesn't see
		// our peering connection yet. Return an empty state.
		if vpcPeeringConnection == nil || vpcPeeringConnection.Status == nil {
			return nil, vpcPeeringConnectionStatusNotFound, nil
		}

		statusCode := aws.StringValue(vpcPeeringConnection.Status.Code)

		// https://docs.aws.amazon.com/vpc/latest/peering/vpc-peering-basics.html#vpc-peering-lifecycle
		switch statusCode {
		case ec2.VpcPeeringConnectionStateReasonCodeFailed:
			log.Printf("[WARN] VPC Peering Connection (%s): %s: %s", vpcPeeringConnectionID, statusCode, aws.StringValue(vpcPeeringConnection.Status.Message))
			fallthrough
		case ec2.VpcPeeringConnectionStateReasonCodeDeleted, ec2.VpcPeeringConnectionStateReasonCodeExpired, ec2.VpcPeeringConnectionStateReasonCodeRejected:
			return nil, vpcPeeringConnectionStatusNotFound, nil
		}

		return vpcPeeringConnection, statusCode, nil
	}
}

const (
	attachmentStateNotFound = "NotFound"
	attachmentStateUnknown  = "Unknown"
)

// statusVPNGatewayVPCAttachmentState fetches the attachment between the specified VPN gateway and VPC and its state
func statusVPNGatewayVPCAttachmentState(conn *ec2.EC2, vpnGatewayID, vpcID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		vpcAttachment, err := findVPNGatewayVPCAttachment(conn, vpnGatewayID, vpcID)
		if tfawserr.ErrCodeEquals(err, invalidVPNGatewayIDNotFound) {
			return nil, attachmentStateNotFound, nil
		}
		if err != nil {
			return nil, attachmentStateUnknown, err
		}

		if vpcAttachment == nil {
			return nil, attachmentStateNotFound, nil
		}

		return vpcAttachment, aws.StringValue(vpcAttachment.State), nil
	}
}

const (
	managedPrefixListStateNotFound = "NotFound"
	managedPrefixListStateUnknown  = "Unknown"
)

func statusManagedPrefixListState(conn *ec2.EC2, prefixListId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		managedPrefixList, err := findManagedPrefixListByID(conn, prefixListId)
		if err != nil {
			return nil, managedPrefixListStateUnknown, err
		}
		if managedPrefixList == nil {
			return nil, managedPrefixListStateNotFound, nil
		}

		return managedPrefixList, aws.StringValue(managedPrefixList.State), nil
	}
}

func statusVPCEndpointState(conn *ec2.EC2, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		vpcEndpoint, err := findVPCEndpointByID(conn, id)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return vpcEndpoint, aws.StringValue(vpcEndpoint.State), nil
	}
}

const (
	vpcEndpointRouteTableAssociationStatusReady = "ready"
)

func statusVPCEndpointRouteTableAssociation(conn *ec2.EC2, vpcEndpointID, routeTableID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		err := findVPCEndpointRouteTableAssociationExists(conn, vpcEndpointID, routeTableID)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return "", vpcEndpointRouteTableAssociationStatusReady, nil
	}
}

func statusEBSSnapshotImport(conn *ec2.EC2, importTaskId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		params := &ec2.DescribeImportSnapshotTasksInput{
			ImportTaskIds: []*string{aws.String(importTaskId)},
		}

		resp, err := conn.DescribeImportSnapshotTasks(params)
		if err != nil {
			return nil, "", err
		}

		if resp == nil || len(resp.ImportSnapshotTasks) < 1 {
			return nil, snapshotImportNotFound, nil
		}

		if task := resp.ImportSnapshotTasks[0]; task != nil {
			detail := task.SnapshotTaskDetail
			if detail.Status != nil && aws.StringValue(detail.Status) == eBSSnapshotImportDeleting {
				err = fmt.Errorf("Snapshot import task is deleting")
			}
			return detail, aws.StringValue(detail.Status), err
		} else {
			return nil, snapshotImportNotFound, nil
		}
	}
}
