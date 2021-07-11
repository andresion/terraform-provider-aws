package provider

import (
	"context"
	"fmt"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/service/amplify"
)

// ServicePackage is the core interface that all service packages must implement.
type ServicePackage interface {
	// Configure is called during provider configuration.
	Configure(context.Context) error

	// DataSources returns a map of the data sources implemented in this service package.
	DataSources() map[string]*schema.Resource

	// ID returns the service package's ID.
	// This ID must be unique.
	ID() string

	// Resources returns a map of the resources implemented in this service package.
	Resources() map[string]*schema.Resource
}

var (
	servicePackageRegistry      map[string]ServicePackage
	servicePackageRegistryError error
	servicePackageRegistryOnce  sync.Once
)

// ServicePackages returns the registered service packages.
func ServicePackages() (map[string]ServicePackage, error) {
	servicePackageRegistryOnce.Do(func() {
		// All the service packages we support.
		servicePackages := []ServicePackage{
			&amplify.ServicePackage{},
		}

		registry := make(map[string]ServicePackage)

		for _, servicePackage := range servicePackages {
			id := servicePackage.ID()

			if _, exists := registry[id]; exists {
				servicePackageRegistryError = fmt.Errorf("A service package with ID %q is already registered", id)

				return
			}

			registry[id] = servicePackage
		}

		servicePackageRegistry = registry
	})

	return servicePackageRegistry, servicePackageRegistryError
}
