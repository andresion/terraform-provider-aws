package registry

import (
	"fmt"
	"sync"

	"github.com/terraform-providers/terraform-provider-aws/aws/internal/provider"
)

var servicePackageRegistrationClosed bool
var servicePackageRegistry map[string]provider.ServicePackage
var servicePackageRegistryMu sync.Mutex

// AddServicePackage registers the specified service package.
func AddServicePackage(servicePackage provider.ServicePackage) error {
	servicePackageRegistryMu.Lock()
	defer servicePackageRegistryMu.Unlock()

	if servicePackageRegistrationClosed {
		return fmt.Errorf("Service package registration is closed")
	}

	if servicePackageRegistry == nil {
		servicePackageRegistry = make(map[string]provider.ServicePackage)
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
func ServicePackages() map[string]provider.ServicePackage {
	servicePackageRegistryMu.Lock()
	defer servicePackageRegistryMu.Unlock()

	servicePackageRegistrationClosed = true

	return servicePackageRegistry
}
