package service

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ServicePackage is the core interface that all service packages must implement.
type ServicePackage interface {
	// Client returns the service's AWS API client.
	Client() interface{}

	// Configure is called during provider configuration.
	Configure(context.Context, client.ConfigProvider) error

	// DataSources returns a map of the data sources implemented in this service package.
	DataSources() map[string]*schema.Resource

	// The ID to lookup a service endpoint with.
	// This is usually (but not always) the AWS SDK EndpointsID constant.
	// See `testAccCheckAWSProviderEndpoints` for exceptions.
	EndpointsID() string

	// ID returns the service package's ID.
	// This ID must be unique.
	ID() string

	// Resources returns a map of the resources implemented in this service package.
	Resources() map[string]*schema.Resource
}
