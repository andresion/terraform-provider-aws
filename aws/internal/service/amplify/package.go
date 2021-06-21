package amplify

import (
	"context"

	"github.com/aws/aws-sdk-go/service/amplify"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/registry"
)

const (
	servicePackageName = amplify.ServiceID
)

func init() {
	if err := registry.AddServicePackage(&servicePackage{}); err != nil {
		panic(err)
	}
}

type servicePackage struct{}

func (sp *servicePackage) Name() string {
	return servicePackageName
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
