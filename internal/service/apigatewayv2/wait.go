package apigatewayv2

import (
	"time"

	"github.com/aws/aws-sdk-go/service/apigatewayv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	// Maximum amount of time to wait for a Deployment to return Deployed
	deploymentDeployedTimeout = 5 * time.Minute

	// Maximum amount of time to wait for a VPC Link to return Available
	vpcLinkAvailableTimeout = 10 * time.Minute

	// Maximum amount of time to wait for a VPC Link to return Deleted
	vpcLinkDeletedTimeout = 10 * time.Minute
)

// waitDeploymentDeployed waits for a Deployment to return Deployed
func waitDeploymentDeployed(conn *apigatewayv2.ApiGatewayV2, apiId, deploymentId string) (*apigatewayv2.GetDeploymentOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{apigatewayv2.DeploymentStatusPending},
		Target:  []string{apigatewayv2.DeploymentStatusDeployed},
		Refresh: statusDeployment(conn, apiId, deploymentId),
		Timeout: deploymentDeployedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if v, ok := outputRaw.(*apigatewayv2.GetDeploymentOutput); ok {
		return v, err
	}

	return nil, err
}

func waitDomainNameAvailable(conn *apigatewayv2.ApiGatewayV2, name string, timeout time.Duration) (*apigatewayv2.GetDomainNameOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{apigatewayv2.DomainNameStatusUpdating},
		Target:  []string{apigatewayv2.DomainNameStatusAvailable},
		Refresh: statusDomainName(conn, name),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if v, ok := outputRaw.(*apigatewayv2.GetDomainNameOutput); ok {
		return v, err
	}

	return nil, err
}

// waitVPCLinkAvailable waits for a VPC Link to return Available
func waitVPCLinkAvailable(conn *apigatewayv2.ApiGatewayV2, vpcLinkId string) (*apigatewayv2.GetVpcLinkOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{apigatewayv2.VpcLinkStatusPending},
		Target:  []string{apigatewayv2.VpcLinkStatusAvailable},
		Refresh: statusVPCLink(conn, vpcLinkId),
		Timeout: vpcLinkAvailableTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if v, ok := outputRaw.(*apigatewayv2.GetVpcLinkOutput); ok {
		return v, err
	}

	return nil, err
}

// waitVPCLinkAvailable waits for a VPC Link to return Deleted
func waitVPCLinkDeleted(conn *apigatewayv2.ApiGatewayV2, vpcLinkId string) (*apigatewayv2.GetVpcLinkOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{apigatewayv2.VpcLinkStatusDeleting},
		Target:  []string{apigatewayv2.VpcLinkStatusFailed},
		Refresh: statusVPCLink(conn, vpcLinkId),
		Timeout: vpcLinkDeletedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if v, ok := outputRaw.(*apigatewayv2.GetVpcLinkOutput); ok {
		return v, err
	}

	return nil, err
}
