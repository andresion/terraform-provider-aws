package elasticache

import (
	"time"

	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	replicationGroupDefaultCreatedTimeout = 60 * time.Minute
	replicationGroupDefaultUpdatedTimeout = 40 * time.Minute
	replicationGroupDefaultDeletedTimeout = 40 * time.Minute

	replicationGroupAvailableMinTimeout = 10 * time.Second
	replicationGroupAvailableDelay      = 30 * time.Second

	replicationGroupDeletedMinTimeout = 10 * time.Second
	replicationGroupDeletedDelay      = 30 * time.Second

	userActiveTimeout  = 5 * time.Minute
	userDeletedTimeout = 5 * time.Minute
)

// waitReplicationGroupAvailable waits for a ReplicationGroup to return Available
func waitReplicationGroupAvailable(conn *elasticache.ElastiCache, replicationGroupID string, timeout time.Duration) (*elasticache.ReplicationGroup, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			replicationGroupStatusCreating,
			replicationGroupStatusModifying,
			replicationGroupStatusSnapshotting,
		},
		Target:     []string{replicationGroupStatusAvailable},
		Refresh:    statusReplicationGroup(conn, replicationGroupID),
		Timeout:    timeout,
		MinTimeout: replicationGroupAvailableMinTimeout,
		Delay:      replicationGroupAvailableDelay,
	}

	outputRaw, err := stateConf.WaitForState()
	if v, ok := outputRaw.(*elasticache.ReplicationGroup); ok {
		return v, err
	}
	return nil, err
}

// waitReplicationGroupDeleted waits for a ReplicationGroup to be deleted
func waitReplicationGroupDeleted(conn *elasticache.ElastiCache, replicationGroupID string, timeout time.Duration) (*elasticache.ReplicationGroup, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			replicationGroupStatusCreating,
			replicationGroupStatusAvailable,
			replicationGroupStatusDeleting,
		},
		Target:     []string{},
		Refresh:    statusReplicationGroup(conn, replicationGroupID),
		Timeout:    timeout,
		MinTimeout: replicationGroupDeletedMinTimeout,
		Delay:      replicationGroupDeletedDelay,
	}

	outputRaw, err := stateConf.WaitForState()
	if v, ok := outputRaw.(*elasticache.ReplicationGroup); ok {
		return v, err
	}
	return nil, err
}

// waitReplicationGroupMemberClustersAvailable waits for all of a ReplicationGroup's Member Clusters to return Available
func waitReplicationGroupMemberClustersAvailable(conn *elasticache.ElastiCache, replicationGroupID string, timeout time.Duration) ([]*elasticache.CacheCluster, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			cacheClusterStatusCreating,
			cacheClusterStatusDeleting,
			cacheClusterStatusModifying,
		},
		Target:     []string{cacheClusterStatusAvailable},
		Refresh:    statusReplicationGroupMemberClusters(conn, replicationGroupID),
		Timeout:    timeout,
		MinTimeout: cacheClusterAvailableMinTimeout,
		Delay:      cacheClusterAvailableDelay,
	}

	outputRaw, err := stateConf.WaitForState()
	if v, ok := outputRaw.([]*elasticache.CacheCluster); ok {
		return v, err
	}
	return nil, err
}

const (
	cacheClusterCreatedTimeout = 40 * time.Minute
	cacheClusterUpdatedTimeout = 80 * time.Minute
	cacheClusterDeletedTimeout = 40 * time.Minute

	cacheClusterAvailableMinTimeout = 10 * time.Second
	cacheClusterAvailableDelay      = 30 * time.Second

	cacheClusterDeletedMinTimeout = 10 * time.Second
	cacheClusterDeletedDelay      = 30 * time.Second
)

// waitCacheClusterAvailable waits for a Cache Cluster to return Available
func waitCacheClusterAvailable(conn *elasticache.ElastiCache, cacheClusterID string, timeout time.Duration) (*elasticache.CacheCluster, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			cacheClusterStatusCreating,
			cacheClusterStatusModifying,
			cacheClusterStatusSnapshotting,
			cacheClusterStatusRebootingClusterNodes,
		},
		Target:     []string{cacheClusterStatusAvailable},
		Refresh:    statusCacheCluster(conn, cacheClusterID),
		Timeout:    timeout,
		MinTimeout: cacheClusterAvailableMinTimeout,
		Delay:      cacheClusterAvailableDelay,
	}

	outputRaw, err := stateConf.WaitForState()
	if v, ok := outputRaw.(*elasticache.CacheCluster); ok {
		return v, err
	}
	return nil, err
}

// waitCacheClusterDeleted waits for a Cache Cluster to be deleted
func waitCacheClusterDeleted(conn *elasticache.ElastiCache, cacheClusterID string, timeout time.Duration) (*elasticache.CacheCluster, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			cacheClusterStatusCreating,
			cacheClusterStatusAvailable,
			cacheClusterStatusModifying,
			cacheClusterStatusDeleting,
			cacheClusterStatusIncompatibleNetwork,
			cacheClusterStatusRestoreFailed,
			cacheClusterStatusSnapshotting,
		},
		Target:     []string{},
		Refresh:    statusCacheCluster(conn, cacheClusterID),
		Timeout:    timeout,
		MinTimeout: cacheClusterDeletedMinTimeout,
		Delay:      cacheClusterDeletedDelay,
	}

	outputRaw, err := stateConf.WaitForState()
	if v, ok := outputRaw.(*elasticache.CacheCluster); ok {
		return v, err
	}
	return nil, err
}

const (
	globalReplicationGroupDefaultCreatedTimeout = 20 * time.Minute
	globalReplicationGroupDefaultUpdatedTimeout = 40 * time.Minute
	globalReplicationGroupDefaultDeletedTimeout = 20 * time.Minute

	globalReplicationGroupAvailableMinTimeout = 10 * time.Second
	globalReplicationGroupAvailableDelay      = 30 * time.Second

	globalReplicationGroupDeletedMinTimeout = 10 * time.Second
	globalReplicationGroupDeletedDelay      = 30 * time.Second
)

// waitGlobalReplicationGroupAvailable waits for a Global Replication Group to be available,
// with status either "available" or "primary-only"
func waitGlobalReplicationGroupAvailable(conn *elasticache.ElastiCache, globalReplicationGroupID string, timeout time.Duration) (*elasticache.GlobalReplicationGroup, error) {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{globalReplicationGroupStatusCreating, globalReplicationGroupStatusModifying},
		Target:     []string{globalReplicationGroupStatusAvailable, globalReplicationGroupStatusPrimaryOnly},
		Refresh:    statusGlobalReplicationGroup(conn, globalReplicationGroupID),
		Timeout:    timeout,
		MinTimeout: globalReplicationGroupAvailableMinTimeout,
		Delay:      globalReplicationGroupAvailableDelay,
	}

	outputRaw, err := stateConf.WaitForState()
	if v, ok := outputRaw.(*elasticache.GlobalReplicationGroup); ok {
		return v, err
	}
	return nil, err
}

// waitGlobalReplicationGroupDeleted waits for a Global Replication Group to be deleted
func waitGlobalReplicationGroupDeleted(conn *elasticache.ElastiCache, globalReplicationGroupID string) (*elasticache.GlobalReplicationGroup, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			globalReplicationGroupStatusAvailable,
			globalReplicationGroupStatusPrimaryOnly,
			globalReplicationGroupStatusModifying,
			globalReplicationGroupStatusDeleting,
		},
		Target:     []string{},
		Refresh:    statusGlobalReplicationGroup(conn, globalReplicationGroupID),
		Timeout:    globalReplicationGroupDefaultDeletedTimeout,
		MinTimeout: globalReplicationGroupDeletedMinTimeout,
		Delay:      globalReplicationGroupDeletedDelay,
	}

	outputRaw, err := stateConf.WaitForState()
	if v, ok := outputRaw.(*elasticache.GlobalReplicationGroup); ok {
		return v, err
	}
	return nil, err
}

const (
	// globalReplicationGroupDisassociationReadyTimeout specifies how long to wait for a global replication group
	// to be in a valid state before disassociating
	globalReplicationGroupDisassociationReadyTimeout = 45 * time.Minute

	// globalReplicationGroupDisassociationTimeout specifies how long to wait for the actual disassociation
	globalReplicationGroupDisassociationTimeout = 20 * time.Minute

	globalReplicationGroupDisassociationMinTimeout = 10 * time.Second
	globalReplicationGroupDisassociationDelay      = 30 * time.Second
)

func waitGlobalReplicationGroupMemberDetached(conn *elasticache.ElastiCache, globalReplicationGroupID, id string) (*elasticache.GlobalReplicationGroupMember, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			globalReplicationGroupMemberStatusAssociated,
		},
		Target:     []string{},
		Refresh:    statusGlobalReplicationGroupMember(conn, globalReplicationGroupID, id),
		Timeout:    globalReplicationGroupDisassociationTimeout,
		MinTimeout: globalReplicationGroupDisassociationMinTimeout,
		Delay:      globalReplicationGroupDisassociationDelay,
	}

	outputRaw, err := stateConf.WaitForState()
	if v, ok := outputRaw.(*elasticache.GlobalReplicationGroupMember); ok {
		return v, err
	}
	return nil, err
}

// waitUserActive waits for an ElastiCache user to reach an active state after modifications
func waitUserActive(conn *elasticache.ElastiCache, userId string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{userStatusModifying},
		Target:  []string{userStatusActive},
		Refresh: statusUser(conn, userId),
		Timeout: userActiveTimeout,
	}

	_, err := stateConf.WaitForState()

	return err
}

// waitUserDeleted waits for an ElastiCache user to be deleted
func waitUserDeleted(conn *elasticache.ElastiCache, userId string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{userStatusDeleting},
		Target:  []string{},
		Refresh: statusUser(conn, userId),
		Timeout: userDeletedTimeout,
	}

	_, err := stateConf.WaitForState()

	return err
}
