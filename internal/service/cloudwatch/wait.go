package cloudwatch

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	metricStreamDeleteTimeout = 2 * time.Minute
	metricStreamReadyTimeout  = 1 * time.Minute

	stateRunning = "running"
	stateStopped = "stopped"
)

func waitMetricStreamDeleted(ctx context.Context, conn *cloudwatch.CloudWatch, name string) (*cloudwatch.GetMetricStreamOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			stateRunning,
			stateStopped,
		},
		Target:  []string{},
		Refresh: statusMetricStreamState(ctx, conn, name),
		Timeout: metricStreamDeleteTimeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if v, ok := outputRaw.(*cloudwatch.GetMetricStreamOutput); ok {
		return v, err
	}

	return nil, err
}

func waitMetricStreamReady(ctx context.Context, conn *cloudwatch.CloudWatch, name string) (*cloudwatch.GetMetricStreamOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			stateStopped,
		},
		Target: []string{
			stateRunning,
		},
		Refresh: statusMetricStreamState(ctx, conn, name),
		Timeout: metricStreamReadyTimeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if v, ok := outputRaw.(*cloudwatch.GetMetricStreamOutput); ok {
		return v, err
	}

	return nil, err
}
