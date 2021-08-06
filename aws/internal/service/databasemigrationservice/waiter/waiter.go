package waiter

import (
	"time"

	dms "github.com/aws/aws-sdk-go/service/databasemigrationservice"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	ReplicationTaskDeletedTimeout = 20 * time.Minute
	ReplicationTaskReadyTimeout   = 20 * time.Minute
	ReplicationTaskStoppedTimeout = 20 * time.Minute
)

func ReplicationTaskDeleted(conn *dms.DatabaseMigrationService, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{ReplicationTaskStatusDeleting},
		Target:     []string{},
		Refresh:    ReplicationTaskStatus(conn, id),
		Timeout:    ReplicationTaskDeletedTimeout,
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second,
	}

	_, err := stateConf.WaitForState()

	return err
}

func ReplicationTaskReady(conn *dms.DatabaseMigrationService, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{ReplicationTaskStatusCreating, ReplicationTaskStatusModifying},
		Target:     []string{ReplicationTaskStatusReady},
		Refresh:    ReplicationTaskStatus(conn, id),
		Timeout:    ReplicationTaskReadyTimeout,
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second,
	}

	_, err := stateConf.WaitForState()

	return err
}

func ReplicationTaskStopped(conn *dms.DatabaseMigrationService, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{ReplicationTaskStatusReady, ReplicationTaskStatusStopping},
		Target:     []string{ReplicationTaskStatusStopped},
		Refresh:    ReplicationTaskStatus(conn, id),
		Timeout:    ReplicationTaskStoppedTimeout,
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second,
	}

	_, err := stateConf.WaitForState()

	return err
}
