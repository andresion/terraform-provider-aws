package ec2

const (
	// https://docs.aws.amazon.com/vpc/latest/privatelink/vpce-interface.html#vpce-interface-lifecycle
	vpcEndpointStateAvailable         = "available"
	vpcEndpointStateDeleted           = "deleted"
	vpcEndpointStateDeleting          = "deleting"
	vpcEndpointStateFailed            = "failed"
	vpcEndpointStatePending           = "pending"
	vpcEndpointStatePendingAcceptance = "pendingAcceptance"
	vpcEndpointStateRejected          = "rejected"
)
