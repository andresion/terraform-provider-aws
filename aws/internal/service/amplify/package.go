package amplify

import (
	"context"

	"github.com/aws/aws-sdk-go/service/amplify"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/provider/meta"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/service"
)

// Implements the service.ServicePackage interface,
type servicePackage struct {
	conn *amplify.Amplify
}

func (sp *servicePackage) ID() string {
	return amplify.ServiceID
}

func (sp *servicePackage) DataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

func (sp *servicePackage) Resources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"aws_amplify_app":                 resourceAwsAmplifyApp(),
		"aws_amplify_backend_environment": resourceAwsAmplifyBackendEnvironment(),
		"aws_amplify_branch":              resourceAwsAmplifyBranch(),
		"aws_amplify_domain_association":  resourceAwsAmplifyDomainAssociation(),
		"aws_amplify_webhook":             resourceAwsAmplifyWebhook(),
	}
}

func (sp *servicePackage) Configure(ctx context.Context) error {
	// TODO Initialize conn.
	return nil
}

func NewServicePackage() service.ServicePackage {
	return &servicePackage{}
}

func fromMeta(v interface{}) (*amplify.Amplify, *keyvaluetags.DefaultConfig, *keyvaluetags.IgnoreConfig) {
	m := v.(meta.Meta)

	return m.GetServicePackage(amplify.ServiceID).(*servicePackage).conn, m.GetDefaultTagsConfig(), m.GetIgnoreTagsConfig()
}
