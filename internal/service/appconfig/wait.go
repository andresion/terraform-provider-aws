package appconfig

import (
	"time"

	"github.com/aws/aws-sdk-go/service/appconfig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const deploymentCreatedTimeout = 20 * time.Minute

func waitDeploymentCreated(conn *appconfig.AppConfig, appID, envID string, deployNum int64) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{appconfig.DeploymentStateBaking, appconfig.DeploymentStateRollingBack, appconfig.DeploymentStateValidating, appconfig.DeploymentStateDeploying},
		Target:  []string{appconfig.DeploymentStateComplete},
		Refresh: statusDeployment(conn, appID, envID, deployNum),
		Timeout: deploymentCreatedTimeout,
	}

	_, err := stateConf.WaitForState()

	return err
}
