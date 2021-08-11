package workspaces

import (
	"time"

	"github.com/aws/aws-sdk-go/service/workspaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	directoryDeregisterInvalidResourceStateTimeout = 2 * time.Minute
	directoryRegisterInvalidResourceStateTimeout   = 2 * time.Minute

	// Maximum amount of time to wait for a Directory to return Registered
	directoryRegisteredTimeout = 10 * time.Minute

	// Maximum amount of time to wait for a Directory to return Deregistered
	directoryDeregisteredTimeout = 10 * time.Minute

	// Maximum amount of time to wait for a WorkSpace to return Available
	workspaceAvailableTimeout = 30 * time.Minute

	// Maximum amount of time to wait for a WorkSpace while returning Updating
	workspaceUpdatingTimeout = 10 * time.Minute

	// Amount of time to delay before checking WorkSpace when updating
	workspaceUpdatingDelay = 1 * time.Minute

	// Maximum amount of time to wait for a WorkSpace to return Terminated
	workspaceTerminatedTimeout = 10 * time.Minute
)

func waitDirectoryRegistered(conn *workspaces.WorkSpaces, directoryID string) (*workspaces.WorkspaceDirectory, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{workspaces.WorkspaceDirectoryStateRegistering},
		Target:  []string{workspaces.WorkspaceDirectoryStateRegistered},
		Refresh: statusDirectoryState(conn, directoryID),
		Timeout: directoryRegisteredTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if v, ok := outputRaw.(*workspaces.WorkspaceDirectory); ok {
		return v, err
	}

	return nil, err
}

func waitDirectoryDeregistered(conn *workspaces.WorkSpaces, directoryID string) (*workspaces.WorkspaceDirectory, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			workspaces.WorkspaceDirectoryStateRegistering,
			workspaces.WorkspaceDirectoryStateRegistered,
			workspaces.WorkspaceDirectoryStateDeregistering,
		},
		Target:  []string{},
		Refresh: statusDirectoryState(conn, directoryID),
		Timeout: directoryDeregisteredTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if v, ok := outputRaw.(*workspaces.WorkspaceDirectory); ok {
		return v, err
	}

	return nil, err
}

func waitWorkspaceAvailable(conn *workspaces.WorkSpaces, workspaceID string, timeout time.Duration) (*workspaces.Workspace, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			workspaces.WorkspaceStatePending,
			workspaces.WorkspaceStateStarting,
		},
		Target:  []string{workspaces.WorkspaceStateAvailable},
		Refresh: statusWorkspaceState(conn, workspaceID),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if v, ok := outputRaw.(*workspaces.Workspace); ok {
		return v, err
	}

	return nil, err
}

func waitWorkspaceTerminated(conn *workspaces.WorkSpaces, workspaceID string, timeout time.Duration) (*workspaces.Workspace, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			workspaces.WorkspaceStatePending,
			workspaces.WorkspaceStateAvailable,
			workspaces.WorkspaceStateImpaired,
			workspaces.WorkspaceStateUnhealthy,
			workspaces.WorkspaceStateRebooting,
			workspaces.WorkspaceStateStarting,
			workspaces.WorkspaceStateRebuilding,
			workspaces.WorkspaceStateRestoring,
			workspaces.WorkspaceStateMaintenance,
			workspaces.WorkspaceStateAdminMaintenance,
			workspaces.WorkspaceStateSuspended,
			workspaces.WorkspaceStateUpdating,
			workspaces.WorkspaceStateStopping,
			workspaces.WorkspaceStateStopped,
			workspaces.WorkspaceStateTerminating,
			workspaces.WorkspaceStateError,
		},
		Target:  []string{workspaces.WorkspaceStateTerminated},
		Refresh: statusWorkspaceState(conn, workspaceID),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if v, ok := outputRaw.(*workspaces.Workspace); ok {
		return v, err
	}

	return nil, err
}

func waitWorkspaceUpdated(conn *workspaces.WorkSpaces, workspaceID string, timeout time.Duration) (*workspaces.Workspace, error) {
	// OperationInProgressException: The properties of this WorkSpace are currently under modification. Please try again in a moment.
	// AWS Workspaces service doesn't change instance status to "Updating" during property modification. Respective AWS Support feature request has been created. Meanwhile, artificial delay is placed here as a workaround.
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			workspaces.WorkspaceStateUpdating,
		},
		Target: []string{
			workspaces.WorkspaceStateAvailable,
			workspaces.WorkspaceStateStopped,
		},
		Refresh: statusWorkspaceState(conn, workspaceID),
		Delay:   workspaceUpdatingDelay,
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if v, ok := outputRaw.(*workspaces.Workspace); ok {
		return v, err
	}

	return nil, err
}
