package finder

import (
	"github.com/aws/aws-sdk-go/aws"
	dms "github.com/aws/aws-sdk-go/service/databasemigrationservice"
)

func ReplicationTask(conn *dms.DatabaseMigrationService, id string) (*dms.ReplicationTask, error) {
	input := &dms.DescribeReplicationTasksInput{
		Filters: []*dms.Filter{
			{
				Name:   aws.String("replication-task-id"),
				Values: []*string{aws.String(id)},
			},
		},
	}

	var result *dms.ReplicationTask

	err := conn.DescribeReplicationTasksPages(input, func(page *dms.DescribeReplicationTasksOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, task := range page.ReplicationTasks {
			if task == nil {
				continue
			}

			if aws.StringValue(task.ReplicationTaskIdentifier) == id {
				result = task
				return false
			}
		}

		return !lastPage
	})

	return result, err
}
