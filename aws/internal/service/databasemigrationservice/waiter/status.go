package waiter

import (
	"github.com/aws/aws-sdk-go/aws"
	dms "github.com/aws/aws-sdk-go/service/databasemigrationservice"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/service/databasemigrationservice/finder"
)

const (
	ReplicationTaskStatusCreating  = "creating"
	ReplicationTaskStatusDeleting  = "deleting"
	ReplicationTaskStatusModifying = "modifying"
	ReplicationTaskStatusReady     = "ready"
	ReplicationTaskStatusRunning   = "running"
	ReplicationTaskStatusStopped   = "stopped"
	ReplicationTaskStatusStopping  = "stopping"
)

func ReplicationTaskStatus(conn *dms.DatabaseMigrationService, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		task, err := finder.ReplicationTask(conn, id)

		if err != nil {
			return nil, "", err
		}

		if task == nil {
			return nil, "", nil
		}

		return task, aws.StringValue(task.Status), nil
	}
}
