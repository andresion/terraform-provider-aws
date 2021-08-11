package mq

import (
	"time"

	"github.com/aws/aws-sdk-go/service/mq"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	brokerCreateTimeout = 30 * time.Minute
	brokerDeleteTimeout = 30 * time.Minute
	brokerRebootTimeout = 30 * time.Minute
)

func waitBrokerCreated(conn *mq.MQ, id string) (*mq.DescribeBrokerResponse, error) {
	stateConf := resource.StateChangeConf{
		Pending: []string{
			mq.BrokerStateCreationInProgress,
			mq.BrokerStateRebootInProgress,
		},
		Target:  []string{mq.BrokerStateRunning},
		Timeout: brokerCreateTimeout,
		Refresh: statusBroker(conn, id),
	}
	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*mq.DescribeBrokerResponse); ok {
		return output, err
	}

	return nil, err
}

func waitBrokerDeleted(conn *mq.MQ, id string) (*mq.DescribeBrokerResponse, error) {
	stateConf := resource.StateChangeConf{
		Pending: []string{
			mq.BrokerStateCreationFailed,
			mq.BrokerStateDeletionInProgress,
			mq.BrokerStateRebootInProgress,
			mq.BrokerStateRunning,
		},
		Target:  []string{},
		Timeout: brokerDeleteTimeout,
		Refresh: statusBroker(conn, id),
	}
	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*mq.DescribeBrokerResponse); ok {
		return output, err
	}

	return nil, err
}

func waitBrokerRebooted(conn *mq.MQ, id string) (*mq.DescribeBrokerResponse, error) {
	stateConf := resource.StateChangeConf{
		Pending: []string{
			mq.BrokerStateRebootInProgress,
		},
		Target:  []string{mq.BrokerStateRunning},
		Timeout: brokerRebootTimeout,
		Refresh: statusBroker(conn, id),
	}
	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*mq.DescribeBrokerResponse); ok {
		return output, err
	}

	return nil, err
}
