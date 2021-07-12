package meta

import (
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/service"
)

// Meta is the interface implemented by the CRUD handlers' `meta` parameter.
type Meta interface {
	// GetAccountID returns the provider's AWS account ID.
	GetAccountID() string

	// GetDefaultTagsConfig returns the provider's `default_tags` configuration.
	GetDefaultTagsConfig() *keyvaluetags.DefaultConfig

	// GetIgnoreTagsConfig returns the provider's `ignore_tags` configuration.
	GetIgnoreTagsConfig() *keyvaluetags.IgnoreConfig

	// GetRegion returns the provider's AWS region.
	GetRegion() string

	// GetServicePackage returns the ServicePackage for the specified service ID.
	GetServicePackage(id string) service.ServicePackage
}
