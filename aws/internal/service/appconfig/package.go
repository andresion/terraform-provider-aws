package appconfig

import (
	"context"

	"github.com/aws/aws-sdk-go/service/appconfig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/provider"
)

func init() {
	if err := provider.RegisterServicePackage(&servicePackage{}); err != nil {
		panic(err)
	}
}

type servicePackage struct{}

func (sp *servicePackage) Name() string {
	return appconfig.ServiceName
}

func (sp *servicePackage) DataSources(ctx context.Context) (map[string]*schema.Resource, error) {
	return map[string]*schema.Resource{}, nil
}

func (sp *servicePackage) Resources(ctx context.Context) (map[string]*schema.Resource, error) {
	return map[string]*schema.Resource{}, nil
}

func (sp *servicePackage) DocumentationCategories(context.Context) ([]string, error) {
	return []string{}, nil
}

func (sp *servicePackage) Configure(ctx context.Context) error {
	return nil
}
