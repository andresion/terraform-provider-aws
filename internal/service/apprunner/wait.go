package apprunner

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/service/apprunner"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	autoScalingConfigurationCreateTimeout = 2 * time.Minute
	autoScalingConfigurationDeleteTimeout = 2 * time.Minute

	connectionDeleteTimeout = 5 * time.Minute

	customDomainAssociationCreateTimeout = 5 * time.Minute
	customDomainAssociationDeleteTimeout = 5 * time.Minute

	serviceCreateTimeout = 20 * time.Minute
	serviceDeleteTimeout = 20 * time.Minute
	serviceUpdateTimeout = 20 * time.Minute
)

func waitAutoScalingConfigurationActive(ctx context.Context, conn *apprunner.AppRunner, arn string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{},
		Target:  []string{autoScalingConfigurationStatusActive},
		Refresh: statusAutoScalingConfiguration(ctx, conn, arn),
		Timeout: autoScalingConfigurationCreateTimeout,
	}

	_, err := stateConf.WaitForState()

	return err
}

func waitAutoScalingConfigurationInactive(ctx context.Context, conn *apprunner.AppRunner, arn string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{autoScalingConfigurationStatusActive},
		Target:  []string{autoScalingConfigurationStatusInactive},
		Refresh: statusAutoScalingConfiguration(ctx, conn, arn),
		Timeout: autoScalingConfigurationDeleteTimeout,
	}

	_, err := stateConf.WaitForState()

	return err
}

func waitConnectionDeleted(ctx context.Context, conn *apprunner.AppRunner, name string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{apprunner.ConnectionStatusPendingHandshake, apprunner.ConnectionStatusAvailable, apprunner.ConnectionStatusDeleted},
		Target:  []string{},
		Refresh: statusConnection(ctx, conn, name),
		Timeout: connectionDeleteTimeout,
	}

	_, err := stateConf.WaitForState()

	return err
}

func waitCustomDomainAssociationCreated(ctx context.Context, conn *apprunner.AppRunner, domainName, serviceArn string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{customDomainAssociationStatusCreating},
		Target:  []string{customDomainAssociationStatusPendingCertificateDNSValidation},
		Refresh: statusCustomDomain(ctx, conn, domainName, serviceArn),
		Timeout: customDomainAssociationCreateTimeout,
	}

	_, err := stateConf.WaitForState()

	return err
}

func waitCustomDomainAssociationDeleted(ctx context.Context, conn *apprunner.AppRunner, domainName, serviceArn string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{customDomainAssociationStatusActive, customDomainAssociationStatusDeleting},
		Target:  []string{},
		Refresh: statusCustomDomain(ctx, conn, domainName, serviceArn),
		Timeout: customDomainAssociationDeleteTimeout,
	}

	_, err := stateConf.WaitForState()

	return err
}

func waitServiceCreated(ctx context.Context, conn *apprunner.AppRunner, serviceArn string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{apprunner.ServiceStatusOperationInProgress},
		Target:  []string{apprunner.ServiceStatusRunning},
		Refresh: statusService(ctx, conn, serviceArn),
		Timeout: serviceCreateTimeout,
	}

	_, err := stateConf.WaitForState()

	return err
}

func waitServiceUpdated(ctx context.Context, conn *apprunner.AppRunner, serviceArn string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{apprunner.ServiceStatusOperationInProgress},
		Target:  []string{apprunner.ServiceStatusRunning},
		Refresh: statusService(ctx, conn, serviceArn),
		Timeout: serviceUpdateTimeout,
	}

	_, err := stateConf.WaitForState()

	return err
}

func waitServiceDeleted(ctx context.Context, conn *apprunner.AppRunner, serviceArn string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{apprunner.ServiceStatusRunning, apprunner.ServiceStatusOperationInProgress},
		Target:  []string{apprunner.ServiceStatusDeleted},
		Refresh: statusService(ctx, conn, serviceArn),
		Timeout: serviceDeleteTimeout,
	}

	_, err := stateConf.WaitForState()

	return err
}
