package sqs

const (
	errCodeInvalidAction = "InvalidAction"
)

const (
	fifoQueueNameSuffix = ".fifo"
)

const (
	defaultQueueDelaySeconds                  = 0
	defaultQueueKMSDataKeyReusePeriodSeconds  = 300
	defaultQueueMaximumMessageSize            = 262_144 // 256 KiB.
	defaultQueueMessageRetentionPeriod        = 345_600 // 4 days.
	defaultQueueReceiveMessageWaitTimeSeconds = 0
	defaultQueueVisibilityTimeout             = 30
)

const (
	deduplicationScopeMessageGroup = "messageGroup"
	deduplicationScopeQueue        = "queue"
)

func deduplicationScope_Values() []string {
	return []string{
		deduplicationScopeMessageGroup,
		deduplicationScopeQueue,
	}
}

const (
	fifoThroughputLimitPerMessageGroupId = "perMessageGroupId"
	fifoThroughputLimitPerQueue          = "perQueue"
)

func fifoThroughputLimit_Values() []string {
	return []string{
		fifoThroughputLimitPerMessageGroupId,
		fifoThroughputLimitPerQueue,
	}
}
