package neptune

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/neptune"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-aws/internal/tfresource"
)

const (
	// EventSubscription NotFound
	eventSubscriptionStatusNotFound = "NotFound"

	// EventSubscription Unknown
	eventSubscriptionStatusUnknown = "Unknown"

	// Cluster NotFound
	clusterStatusNotFound = "NotFound"

	// Cluster Unknown
	clusterStatusUnknown = "Unknown"

	// DBClusterEndpoint Unknown
	dbClusterEndpointStatusUnknown = "Unknown"
)

// statusEventSubscription fetches the EventSubscription and its Status
func statusEventSubscription(conn *neptune.Neptune, subscriptionName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		input := &neptune.DescribeEventSubscriptionsInput{
			SubscriptionName: aws.String(subscriptionName),
		}

		output, err := conn.DescribeEventSubscriptions(input)

		if err != nil {
			return nil, eventSubscriptionStatusUnknown, err
		}

		if len(output.EventSubscriptionsList) == 0 {
			return nil, eventSubscriptionStatusNotFound, nil
		}

		return output.EventSubscriptionsList[0], aws.StringValue(output.EventSubscriptionsList[0].Status), nil
	}
}

// statusCluster fetches the Cluster and its Status
func statusCluster(conn *neptune.Neptune, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		input := &neptune.DescribeDBClustersInput{
			DBClusterIdentifier: aws.String(id),
		}

		output, err := conn.DescribeDBClusters(input)

		if err != nil {
			return nil, clusterStatusUnknown, err
		}

		if len(output.DBClusters) == 0 {
			return nil, clusterStatusNotFound, nil
		}

		cluster := output.DBClusters[0]

		return cluster, aws.StringValue(cluster.Status), nil
	}
}

// statusDBClusterEndpoint fetches the DBClusterEndpoint and its Status
func statusDBClusterEndpoint(conn *neptune.Neptune, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := findEndpointById(conn, id)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, dbClusterEndpointStatusUnknown, err
		}

		return output, aws.StringValue(output.Status), nil
	}
}
