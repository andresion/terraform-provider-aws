package registry

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
	DataSources() map[string]*schema.Resource

	// Name returns the service package's name.
	// This name must be unique.
	Name() string

	// Resources returns a map of the resources implemented in this service package.
	Resources() map[string]*schema.Resource
}

var servicePackageRegistrationClosed bool
var servicePackageRegistry map[string]ServicePackage
var servicePackageRegistryMu sync.Mutex

// AddServicePackage registers the specified service package.
func AddServicePackage(servicePackage ServicePackage) error {
	servicePackageRegistryMu.Lock()
	defer servicePackageRegistryMu.Unlock()

	if servicePackageRegistrationClosed {
		return fmt.Errorf("Service package registration is closed")
	}

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

// ServicePackages returns the registered service packages.
// Service package registration is closed.
func ServicePackages() map[string]ServicePackage {
	servicePackageRegistryMu.Lock()
	defer servicePackageRegistryMu.Unlock()

	servicePackageRegistrationClosed = true

	return servicePackageRegistry
}
