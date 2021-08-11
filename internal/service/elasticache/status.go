package elasticache

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-aws/internal/tfresource"
)

const (
	replicationGroupStatusCreating     = "creating"
	replicationGroupStatusAvailable    = "available"
	replicationGroupStatusModifying    = "modifying"
	replicationGroupStatusDeleting     = "deleting"
	replicationGroupStatusCreateFailed = "create-failed"
	replicationGroupStatusSnapshotting = "snapshotting"

	userStatusActive    = "active"
	userStatusDeleting  = "deleting"
	userStatusModifying = "modifying"
)

// statusReplicationGroup fetches the Replication Group and its Status
func statusReplicationGroup(conn *elasticache.ElastiCache, replicationGroupID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		rg, err := findReplicationGroupByID(conn, replicationGroupID)
		if tfresource.NotFound(err) {
			return nil, "", nil
		}
		if err != nil {
			return nil, "", err
		}

		return rg, aws.StringValue(rg.Status), nil
	}
}

// statusReplicationGroupMemberClusters fetches the Replication Group's Member Clusters and either "available" or the first non-"available" status.
// NOTE: This function assumes that the intended end-state is to have all member clusters in "available" status.
func statusReplicationGroupMemberClusters(conn *elasticache.ElastiCache, replicationGroupID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		clusters, err := findReplicationGroupMemberClustersByID(conn, replicationGroupID)
		if tfresource.NotFound(err) {
			return nil, "", nil
		}
		if err != nil {
			return nil, "", err
		}

		status := cacheClusterStatusAvailable
		for _, v := range clusters {
			clusterStatus := aws.StringValue(v.CacheClusterStatus)
			if clusterStatus != cacheClusterStatusAvailable {
				status = clusterStatus
				break
			}
		}
		return clusters, status, nil
	}
}

const (
	cacheClusterStatusAvailable             = "available"
	cacheClusterStatusCreating              = "creating"
	cacheClusterStatusDeleted               = "deleted"
	cacheClusterStatusDeleting              = "deleting"
	cacheClusterStatusIncompatibleNetwork   = "incompatible-network"
	cacheClusterStatusModifying             = "modifying"
	cacheClusterStatusRebootingClusterNodes = "rebooting cluster nodes"
	cacheClusterStatusRestoreFailed         = "restore-failed"
	cacheClusterStatusSnapshotting          = "snapshotting"
)

// statusCacheCluster fetches the Cache Cluster and its Status
func statusCacheCluster(conn *elasticache.ElastiCache, cacheClusterID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		c, err := findCacheClusterByID(conn, cacheClusterID)
		if tfresource.NotFound(err) {
			return nil, "", nil
		}
		if err != nil {
			return nil, "", err
		}

		return c, aws.StringValue(c.CacheClusterStatus), nil
	}
}

const (
	globalReplicationGroupStatusAvailable   = "available"
	globalReplicationGroupStatusCreating    = "creating"
	globalReplicationGroupStatusModifying   = "modifying"
	globalReplicationGroupStatusPrimaryOnly = "primary-only"
	globalReplicationGroupStatusDeleting    = "deleting"
	globalReplicationGroupStatusDeleted     = "deleted"
)

// statusGlobalReplicationGroup fetches the Global Replication Group and its Status
func statusGlobalReplicationGroup(conn *elasticache.ElastiCache, globalReplicationGroupID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		grg, err := findGlobalReplicationGroupByID(conn, globalReplicationGroupID)
		if tfresource.NotFound(err) {
			return nil, "", nil
		}
		if err != nil {
			return nil, "", err
		}

		return grg, aws.StringValue(grg.Status), nil
	}
}

const (
	globalReplicationGroupMemberStatusAssociated = "associated"
)

// statusGlobalReplicationGroup fetches the Global Replication Group and its Status
func statusGlobalReplicationGroupMember(conn *elasticache.ElastiCache, globalReplicationGroupID, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		member, err := findGlobalReplicationGroupMemberByID(conn, globalReplicationGroupID, id)
		if tfresource.NotFound(err) {
			return nil, "", nil
		}
		if err != nil {
			return nil, "", err
		}

		return member, aws.StringValue(member.Status), nil
	}
}

// statusUser fetches the ElastiCache user and its Status
func statusUser(conn *elasticache.ElastiCache, userId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		user, err := findElastiCacheUserById(conn, userId)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return user, aws.StringValue(user.Status), nil
	}
}
