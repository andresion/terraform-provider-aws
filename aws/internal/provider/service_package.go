package provider

import (
	"context"
	"fmt"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ServicePackage is the core interface that all service packages must implement.
type ServicePackage interface {
	// Configure is called during provider configuration.
	Configure(context.Context) error

	// DataSources returns a map of the data sources implemented in this service package.
	DataSources(context.Context) (map[string]*schema.Resource, error)

	// DocumentationCategories returns a list of categories which can be used for the documentation sidebar.
	DocumentationCategories(context.Context) ([]string, error)

	// Name returns the service package's name.
	// This name must be unique.
	Name() string

	// Resources returns a map of the resources implemented in this service package.
	Resources(context.Context) (map[string]*schema.Resource, error)
}

var servicePackageRegistry map[string]ServicePackage
var servicePackageRegistryMu sync.Mutex

// RegisterServicePackage registers the specified service package.
func RegisterServicePackage(servicePackage ServicePackage) error {
	servicePackageRegistryMu.Lock()
	defer servicePackageRegistryMu.Unlock()

	if servicePackageRegistry == nil {
		servicePackageRegistry = make(map[string]ServicePackage)
	}

	name := servicePackage.Name()

	if _, exists := servicePackageRegistry[name]; exists {
		return fmt.Errorf("A service package named %q is already registered", name)
	}

	servicePackageRegistry[name] = servicePackage

	return nil
}
