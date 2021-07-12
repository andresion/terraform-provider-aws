package provider

import (
	"fmt"
	"sync"

	"github.com/terraform-providers/terraform-provider-aws/aws/internal/service"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/service/amplify"
)

var (
	servicePackageRegistry      map[string]service.ServicePackage
	servicePackageRegistryError error
	servicePackageRegistryOnce  sync.Once
)

// ServicePackages returns the registered service packages.
func ServicePackages() (map[string]service.ServicePackage, error) {
	servicePackageRegistryOnce.Do(func() {
		// All the service packages we support.
		servicePackages := []service.ServicePackage{
			amplify.NewServicePackage(),
		}

		registry := make(map[string]service.ServicePackage)

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
