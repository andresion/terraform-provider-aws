package servicecatalog

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicecatalog"
	"github.com/hashicorp/aws-sdk-go-base/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-aws/internal/tfresource"
)

const (
	productReadyTimeout  = 3 * time.Minute
	productDeleteTimeout = 3 * time.Minute

	tagOptionReadyTimeout  = 3 * time.Minute
	tagOptionDeleteTimeout = 3 * time.Minute

	portfolioShareCreateTimeout = 3 * time.Minute

	organizationsAccessStableTimeout = 3 * time.Minute
	constraintReadyTimeout           = 3 * time.Minute
	constraintDeleteTimeout          = 3 * time.Minute

	productPortfolioAssociationReadyTimeout  = 3 * time.Minute
	productPortfolioAssociationDeleteTimeout = 3 * time.Minute

	serviceActionReadyTimeout  = 3 * time.Minute
	serviceActionDeleteTimeout = 3 * time.Minute

	budgetResourceAssociationReadyTimeout  = 3 * time.Minute
	budgetResourceAssociationDeleteTimeout = 3 * time.Minute

	tagOptionResourceAssociationReadyTimeout  = 3 * time.Minute
	tagOptionResourceAssociationDeleteTimeout = 3 * time.Minute

	provisioningArtifactReadyTimeout   = 3 * time.Minute
	provisioningArtifactDeletedTimeout = 3 * time.Minute

	principalPortfolioAssociationReadyTimeout  = 3 * time.Minute
	principalPortfolioAssociationDeleteTimeout = 3 * time.Minute

	launchPathsReadyTimeout = 3 * time.Minute

	provisionedProductReadyTimeout  = 30 * time.Minute
	provisionedProductUpdateTimeout = 30 * time.Minute
	provisionedProductDeleteTimeout = 30 * time.Minute

	recordReadyTimeout = 30 * time.Minute

	portfolioConstraintsReadyTimeout = 3 * time.Minute

	minTimeout                 = 2 * time.Second
	notFoundChecks             = 5
	continuousTargetOccurrence = 2

	statusNotFound    = "NOT_FOUND"
	statusUnavailable = "UNAVAILABLE"

	// AWS documentation is wrong, says that status will be "AVAILABLE" but it is actually "CREATED"
	statusCreated = "CREATED"

	organizationAccessStatusError = "ERROR"
)

func waitProductReady(conn *servicecatalog.ServiceCatalog, acceptLanguage, productID string) (*servicecatalog.DescribeProductAsAdminOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{servicecatalog.StatusCreating, statusNotFound, statusUnavailable},
		Target:                    []string{servicecatalog.StatusAvailable, statusCreated},
		Refresh:                   statusProduct(conn, acceptLanguage, productID),
		Timeout:                   productReadyTimeout,
		ContinuousTargetOccurence: continuousTargetOccurrence,
		notFoundChecks:            notFoundChecks,
		minTimeout:                minTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*servicecatalog.DescribeProductAsAdminOutput); ok {
		return output, err
	}

	return nil, err
}

func waitProductDeleted(conn *servicecatalog.ServiceCatalog, acceptLanguage, productID string) (*servicecatalog.DescribeProductAsAdminOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{servicecatalog.StatusCreating, servicecatalog.StatusAvailable, statusCreated, statusUnavailable},
		Target:  []string{statusNotFound},
		Refresh: statusProduct(conn, acceptLanguage, productID),
		Timeout: productDeleteTimeout,
	}

	_, err := stateConf.WaitForState()

	if tfawserr.ErrCodeEquals(err, servicecatalog.ErrCodeResourceNotFoundException) {
		return nil, nil
	}

	return nil, err
}

func waitTagOptionReady(conn *servicecatalog.ServiceCatalog, id string) (*servicecatalog.TagOptionDetail, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{statusNotFound, statusUnavailable},
		Target:  []string{servicecatalog.StatusAvailable},
		Refresh: statusTagOption(conn, id),
		Timeout: tagOptionReadyTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*servicecatalog.TagOptionDetail); ok {
		return output, err
	}

	return nil, err
}

func waitTagOptionDeleted(conn *servicecatalog.ServiceCatalog, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{servicecatalog.StatusAvailable},
		Target:  []string{statusNotFound, statusUnavailable},
		Refresh: statusTagOption(conn, id),
		Timeout: tagOptionDeleteTimeout,
	}

	_, err := stateConf.WaitForState()

	if tfawserr.ErrCodeEquals(err, servicecatalog.ErrCodeResourceNotFoundException) {
		return nil
	}

	return err
}

func waitPortfolioShareReady(conn *servicecatalog.ServiceCatalog, portfolioID, shareType, principalID string, acceptRequired bool) (*servicecatalog.PortfolioShareDetail, error) {
	targets := []string{servicecatalog.ShareStatusCompleted}

	if !acceptRequired {
		targets = append(targets, servicecatalog.ShareStatusInProgress)
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{servicecatalog.ShareStatusNotStarted, servicecatalog.ShareStatusInProgress, statusNotFound, statusUnavailable},
		Target:  targets,
		Refresh: statusPortfolioShare(conn, portfolioID, shareType, principalID),
		Timeout: portfolioShareCreateTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*servicecatalog.PortfolioShareDetail); ok {
		return output, err
	}

	return nil, err
}

func waitPortfolioShareCreatedWithToken(conn *servicecatalog.ServiceCatalog, token string, acceptRequired bool) (*servicecatalog.DescribePortfolioShareStatusOutput, error) {
	targets := []string{servicecatalog.ShareStatusCompleted}

	if !acceptRequired {
		targets = append(targets, servicecatalog.ShareStatusInProgress)
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{servicecatalog.ShareStatusNotStarted, servicecatalog.ShareStatusInProgress, statusNotFound, statusUnavailable},
		Target:  targets,
		Refresh: statusPortfolioShareWithToken(conn, token),
		Timeout: portfolioShareCreateTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*servicecatalog.DescribePortfolioShareStatusOutput); ok {
		return output, err
	}

	return nil, err
}

func waitPortfolioShareDeleted(conn *servicecatalog.ServiceCatalog, portfolioID, shareType, principalID string) (*servicecatalog.PortfolioShareDetail, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{servicecatalog.ShareStatusNotStarted, servicecatalog.ShareStatusInProgress, servicecatalog.ShareStatusCompleted, statusUnavailable},
		Target:  []string{statusNotFound},
		Refresh: statusPortfolioShare(conn, portfolioID, shareType, principalID),
		Timeout: portfolioShareCreateTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if tfresource.NotFound(err) {
		return nil, nil
	}

	if output, ok := outputRaw.(*servicecatalog.PortfolioShareDetail); ok {
		return output, err
	}

	return nil, err
}

func waitPortfolioShareDeletedWithToken(conn *servicecatalog.ServiceCatalog, token string) (*servicecatalog.DescribePortfolioShareStatusOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{servicecatalog.ShareStatusNotStarted, servicecatalog.ShareStatusInProgress, statusNotFound, statusUnavailable},
		Target:  []string{servicecatalog.ShareStatusCompleted},
		Refresh: statusPortfolioShareWithToken(conn, token),
		Timeout: portfolioShareCreateTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*servicecatalog.DescribePortfolioShareStatusOutput); ok {
		return output, err
	}

	return nil, err
}

func waitOrganizationsAccessStable(conn *servicecatalog.ServiceCatalog) (string, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{servicecatalog.AccessStatusUnderChange, statusNotFound, statusUnavailable},
		Target:  []string{servicecatalog.AccessStatusEnabled, servicecatalog.AccessStatusDisabled},
		Refresh: statusOrganizationsAccess(conn),
		Timeout: organizationsAccessStableTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*servicecatalog.GetAWSOrganizationsAccessStatusOutput); ok {
		return aws.StringValue(output.AccessStatus), err
	}

	return "", err
}

func waitConstraintReady(conn *servicecatalog.ServiceCatalog, acceptLanguage, id string) (*servicecatalog.DescribeConstraintOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{statusNotFound, servicecatalog.StatusCreating, statusUnavailable},
		Target:                    []string{servicecatalog.StatusAvailable},
		Refresh:                   statusConstraint(conn, acceptLanguage, id),
		Timeout:                   constraintReadyTimeout,
		ContinuousTargetOccurence: continuousTargetOccurrence,
		notFoundChecks:            notFoundChecks,
		minTimeout:                minTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*servicecatalog.DescribeConstraintOutput); ok {
		return output, err
	}

	return nil, err
}

func waitConstraintDeleted(conn *servicecatalog.ServiceCatalog, acceptLanguage, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{servicecatalog.StatusAvailable, servicecatalog.StatusCreating},
		Target:  []string{statusNotFound},
		Refresh: statusConstraint(conn, acceptLanguage, id),
		Timeout: constraintDeleteTimeout,
	}

	_, err := stateConf.WaitForState()

	return err
}

func waitProductPortfolioAssociationReady(conn *servicecatalog.ServiceCatalog, acceptLanguage, portfolioID, productID string) (*servicecatalog.PortfolioDetail, error) {
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{statusNotFound, statusUnavailable},
		Target:                    []string{servicecatalog.StatusAvailable},
		Refresh:                   statusProductPortfolioAssociation(conn, acceptLanguage, portfolioID, productID),
		Timeout:                   productPortfolioAssociationReadyTimeout,
		ContinuousTargetOccurence: continuousTargetOccurrence,
		notFoundChecks:            notFoundChecks,
		minTimeout:                minTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*servicecatalog.PortfolioDetail); ok {
		return output, err
	}

	return nil, err
}

func waitProductPortfolioAssociationDeleted(conn *servicecatalog.ServiceCatalog, acceptLanguage, portfolioID, productID string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{servicecatalog.StatusAvailable},
		Target:  []string{statusNotFound, statusUnavailable},
		Refresh: statusProductPortfolioAssociation(conn, acceptLanguage, portfolioID, productID),
		Timeout: productPortfolioAssociationDeleteTimeout,
	}

	_, err := stateConf.WaitForState()

	return err
}

func waitServiceActionReady(conn *servicecatalog.ServiceCatalog, acceptLanguage, id string) (*servicecatalog.ServiceActionDetail, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{statusNotFound, statusUnavailable},
		Target:  []string{servicecatalog.StatusAvailable},
		Refresh: statusServiceAction(conn, acceptLanguage, id),
		Timeout: serviceActionReadyTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*servicecatalog.ServiceActionDetail); ok {
		return output, err
	}

	return nil, err
}

func waitServiceActionDeleted(conn *servicecatalog.ServiceCatalog, acceptLanguage, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{servicecatalog.StatusAvailable},
		Target:  []string{statusNotFound, statusUnavailable},
		Refresh: statusServiceAction(conn, acceptLanguage, id),
		Timeout: serviceActionDeleteTimeout,
	}

	_, err := stateConf.WaitForState()

	if tfawserr.ErrCodeEquals(err, servicecatalog.ErrCodeResourceNotFoundException) {
		return nil
	}

	return err
}

func waitBudgetResourceAssociationReady(conn *servicecatalog.ServiceCatalog, budgetName, resourceID string) (*servicecatalog.BudgetDetail, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{statusNotFound, statusUnavailable},
		Target:  []string{servicecatalog.StatusAvailable},
		Refresh: statusBudgetResourceAssociation(conn, budgetName, resourceID),
		Timeout: budgetResourceAssociationReadyTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*servicecatalog.BudgetDetail); ok {
		return output, err
	}

	return nil, err
}

func waitBudgetResourceAssociationDeleted(conn *servicecatalog.ServiceCatalog, budgetName, resourceID string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{servicecatalog.StatusAvailable},
		Target:  []string{statusNotFound, statusUnavailable},
		Refresh: statusBudgetResourceAssociation(conn, budgetName, resourceID),
		Timeout: budgetResourceAssociationDeleteTimeout,
	}

	_, err := stateConf.WaitForState()

	return err
}

func waitTagOptionResourceAssociationReady(conn *servicecatalog.ServiceCatalog, tagOptionID, resourceID string) (*servicecatalog.ResourceDetail, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{statusNotFound, statusUnavailable},
		Target:  []string{servicecatalog.StatusAvailable},
		Refresh: statusTagOptionResourceAssociation(conn, tagOptionID, resourceID),
		Timeout: tagOptionResourceAssociationReadyTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*servicecatalog.ResourceDetail); ok {
		return output, err
	}

	return nil, err
}

func waitTagOptionResourceAssociationDeleted(conn *servicecatalog.ServiceCatalog, tagOptionID, resourceID string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{servicecatalog.StatusAvailable},
		Target:  []string{statusNotFound, statusUnavailable},
		Refresh: statusTagOptionResourceAssociation(conn, tagOptionID, resourceID),
		Timeout: tagOptionResourceAssociationDeleteTimeout,
	}

	_, err := stateConf.WaitForState()

	return err
}

func waitProvisioningArtifactReady(conn *servicecatalog.ServiceCatalog, id, productID string) (*servicecatalog.DescribeProvisioningArtifactOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{servicecatalog.StatusCreating, statusNotFound, statusUnavailable},
		Target:                    []string{servicecatalog.StatusAvailable, statusCreated},
		Refresh:                   statusProvisioningArtifact(conn, id, productID),
		Timeout:                   provisioningArtifactReadyTimeout,
		ContinuousTargetOccurence: continuousTargetOccurrence,
		notFoundChecks:            notFoundChecks,
		minTimeout:                minTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*servicecatalog.DescribeProvisioningArtifactOutput); ok {
		return output, err
	}

	return nil, err
}

func waitProvisioningArtifactDeleted(conn *servicecatalog.ServiceCatalog, id, productID string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{servicecatalog.StatusCreating, servicecatalog.StatusAvailable, statusCreated, statusUnavailable},
		Target:  []string{statusNotFound},
		Refresh: statusProvisioningArtifact(conn, id, productID),
		Timeout: provisioningArtifactDeletedTimeout,
	}

	_, err := stateConf.WaitForState()

	if tfawserr.ErrCodeEquals(err, servicecatalog.ErrCodeResourceNotFoundException) {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}

func waitPrincipalPortfolioAssociationReady(conn *servicecatalog.ServiceCatalog, acceptLanguage, principalARN, portfolioID string) (*servicecatalog.Principal, error) {
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{statusNotFound, statusUnavailable},
		Target:                    []string{servicecatalog.StatusAvailable},
		Refresh:                   statusPrincipalPortfolioAssociation(conn, acceptLanguage, principalARN, portfolioID),
		Timeout:                   principalPortfolioAssociationReadyTimeout,
		ContinuousTargetOccurence: continuousTargetOccurrence,
		notFoundChecks:            notFoundChecks,
		minTimeout:                minTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*servicecatalog.Principal); ok {
		return output, err
	}

	return nil, err
}

func waitPrincipalPortfolioAssociationDeleted(conn *servicecatalog.ServiceCatalog, acceptLanguage, principalARN, portfolioID string) error {
	stateConf := &resource.StateChangeConf{
		Pending:        []string{servicecatalog.StatusAvailable},
		Target:         []string{statusNotFound, statusUnavailable},
		Refresh:        statusPrincipalPortfolioAssociation(conn, acceptLanguage, principalARN, portfolioID),
		Timeout:        principalPortfolioAssociationDeleteTimeout,
		notFoundChecks: 1,
	}

	_, err := stateConf.WaitForState()

	return err
}

func waitLaunchPathsReady(conn *servicecatalog.ServiceCatalog, acceptLanguage, productID string) ([]*servicecatalog.LaunchPathSummary, error) {
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{statusNotFound},
		Target:                    []string{servicecatalog.StatusAvailable},
		Refresh:                   statusLaunchPaths(conn, acceptLanguage, productID),
		Timeout:                   launchPathsReadyTimeout,
		ContinuousTargetOccurence: continuousTargetOccurrence,
		notFoundChecks:            notFoundChecks,
		minTimeout:                minTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.([]*servicecatalog.LaunchPathSummary); ok {
		return output, err
	}

	return nil, err
}

func waitProvisionedProductReady(conn *servicecatalog.ServiceCatalog, acceptLanguage, id, name string) (*servicecatalog.DescribeProvisionedProductOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{statusNotFound, statusUnavailable, servicecatalog.ProvisionedProductStatusUnderChange, servicecatalog.ProvisionedProductStatusPlanInProgress},
		Target:                    []string{servicecatalog.StatusAvailable},
		Refresh:                   statusProvisionedProduct(conn, acceptLanguage, id, name),
		Timeout:                   provisionedProductReadyTimeout,
		ContinuousTargetOccurence: continuousTargetOccurrence,
		notFoundChecks:            notFoundChecks,
		minTimeout:                minTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*servicecatalog.DescribeProvisionedProductOutput); ok {
		return output, err
	}

	return nil, err
}

func waitProvisionedProductTerminated(conn *servicecatalog.ServiceCatalog, acceptLanguage, id, name string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{servicecatalog.StatusAvailable, servicecatalog.ProvisionedProductStatusUnderChange},
		Target:  []string{statusNotFound, statusUnavailable},
		Refresh: statusProvisionedProduct(conn, acceptLanguage, id, name),
		Timeout: provisionedProductDeleteTimeout,
	}

	_, err := stateConf.WaitForState()

	return err
}

func waitRecordReady(conn *servicecatalog.ServiceCatalog, acceptLanguage, id string) (*servicecatalog.DescribeRecordOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{statusNotFound, statusUnavailable, servicecatalog.ProvisionedProductStatusUnderChange, servicecatalog.ProvisionedProductStatusPlanInProgress},
		Target:                    []string{servicecatalog.RecordStatusSucceeded, servicecatalog.StatusAvailable},
		Refresh:                   statusRecord(conn, acceptLanguage, id),
		Timeout:                   recordReadyTimeout,
		ContinuousTargetOccurence: continuousTargetOccurrence,
		notFoundChecks:            notFoundChecks,
		minTimeout:                minTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*servicecatalog.DescribeRecordOutput); ok {
		return output, err
	}

	return nil, err
}

func waitPortfolioConstraintsReady(conn *servicecatalog.ServiceCatalog, acceptLanguage, portfolioID, productID string) ([]*servicecatalog.ConstraintDetail, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{statusNotFound},
		Target:  []string{servicecatalog.StatusAvailable},
		Refresh: statusPortfolioConstraints(conn, acceptLanguage, portfolioID, productID),
		Timeout: portfolioConstraintsReadyTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.([]*servicecatalog.ConstraintDetail); ok {
		return output, err
	}

	return nil, err
}
