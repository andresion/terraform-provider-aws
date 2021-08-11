package neptune

import (
	"time"

	"github.com/aws/aws-sdk-go/service/neptune"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	// Maximum amount of time to wait for an EventSubscription to return Deleted
	eventSubscriptionDeletedTimeout = 10 * time.Minute

	// Maximum amount of time to wait for an DBClusterEndpoint to return Available
	dbClusterEndpointAvailableTimeout = 10 * time.Minute

	// Maximum amount of time to wait for an DBClusterEndpoint to return Deleted
	dbClusterEndpointDeletedTimeout = 10 * time.Minute
)

// waitEventSubscriptionDeleted waits for a EventSubscription to return Deleted
func waitEventSubscriptionDeleted(conn *neptune.Neptune, subscriptionName string) (*neptune.EventSubscription, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"deleting"},
		Target:  []string{eventSubscriptionStatusNotFound},
		Refresh: statusEventSubscription(conn, subscriptionName),
		Timeout: eventSubscriptionDeletedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if v, ok := outputRaw.(*neptune.EventSubscription); ok {
		return v, err
	}

	return nil, err
}

// waitDBClusterDeleted waits for a Cluster to return Deleted
func waitDBClusterDeleted(conn *neptune.Neptune, id string, timeout time.Duration) (*neptune.DBCluster, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			"available",
			"deleting",
			"backing-up",
			"modifying",
		},
		Target:     []string{clusterStatusNotFound},
		Refresh:    statusCluster(conn, id),
		Timeout:    timeout,
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second,
	}

	outputRaw, err := stateConf.WaitForState()

	if v, ok := outputRaw.(*neptune.DBCluster); ok {
		return v, err
	}

	return nil, err
}

// waitDBClusterAvailable waits for a Cluster to return Available
func waitDBClusterAvailable(conn *neptune.Neptune, id string, timeout time.Duration) (*neptune.DBCluster, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			"creating",
			"backing-up",
			"modifying",
			"preparing-data-migration",
			"migrating",
			"configuring-iam-database-auth",
		},
		Target:     []string{"available"},
		Refresh:    statusCluster(conn, id),
		Timeout:    timeout,
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second,
	}

	outputRaw, err := stateConf.WaitForState()

	if v, ok := outputRaw.(*neptune.DBCluster); ok {
		return v, err
	}

	return nil, err
}

// waitDBClusterEndpointAvailable waits for a DBClusterEndpoint to return Available
func waitDBClusterEndpointAvailable(conn *neptune.Neptune, id string) (*neptune.DBClusterEndpoint, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"creating", "modifying"},
		Target:  []string{"available"},
		Refresh: statusDBClusterEndpoint(conn, id),
		Timeout: dbClusterEndpointAvailableTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if v, ok := outputRaw.(*neptune.DBClusterEndpoint); ok {
		return v, err
	}

	return nil, err
}

// waitDBClusterEndpointDeleted waits for a DBClusterEndpoint to return Deleted
func waitDBClusterEndpointDeleted(conn *neptune.Neptune, id string) (*neptune.DBClusterEndpoint, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"deleting"},
		Target:  []string{},
		Refresh: statusDBClusterEndpoint(conn, id),
		Timeout: dbClusterEndpointDeletedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if v, ok := outputRaw.(*neptune.DBClusterEndpoint); ok {
		return v, err
	}

	return nil, err
}
