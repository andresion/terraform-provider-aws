package servicediscovery

import (
	"time"

	"github.com/aws/aws-sdk-go/service/servicediscovery"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	// Maximum amount of time to wait for an Operation to return Success
	operationSuccessTimeout = 5 * time.Minute
)

// waitOperationSuccess waits for an Operation to return Success
func waitOperationSuccess(conn *servicediscovery.ServiceDiscovery, operationID string) (*servicediscovery.GetOperationOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{servicediscovery.OperationStatusSubmitted, servicediscovery.OperationStatusPending},
		Target:  []string{servicediscovery.OperationStatusSuccess},
		Refresh: statusOperation(conn, operationID),
		Timeout: operationSuccessTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*servicediscovery.GetOperationOutput); ok {
		return output, err
	}

	return nil, err
}
