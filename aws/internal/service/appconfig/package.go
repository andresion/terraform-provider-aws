package appconfig

import (
	"context"

	"github.com/aws/aws-sdk-go/service/appconfig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/registry"
)

func init() {
	if err := registry.AddServicePackage(&servicePackage{}); err != nil {
		panic(err)
	}
}

type servicePackage struct{}

func (sp *servicePackage) Name() string {
	return appconfig.ServiceName
}

func (sp *servicePackage) DataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

func (sp *servicePackage) Resources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

func (sp *servicePackage) Configure(ctx context.Context) error {
	return nil
}
