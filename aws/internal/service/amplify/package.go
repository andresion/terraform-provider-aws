package amplify

import (
	"context"

	"github.com/aws/aws-sdk-go/service/amplify"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
)

// Implements the provider.ServicePackage interface,
type ServicePackage struct {
	conn *amplify.Amplify
}

func (sp *ServicePackage) ID() string {
	return amplify.ServiceID
}

func (sp *ServicePackage) DataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

func (sp *ServicePackage) Resources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"aws_amplify_app":                 resourceAwsAmplifyApp(),
		"aws_amplify_backend_environment": resourceAwsAmplifyBackendEnvironment(),
		"aws_amplify_branch":              resourceAwsAmplifyBranch(),
		"aws_amplify_domain_association":  resourceAwsAmplifyDomainAssociation(),
		"aws_amplify_webhook":             resourceAwsAmplifyWebhook(),
	}
}

func (sp *ServicePackage) Configure(ctx context.Context) error {
	// TODO Initialize conn.
	return nil
}

// TODO Consolidate into a single internal package.
type Meta interface {
	GetDefaultTagsConfig() *keyvaluetags.DefaultConfig
	GetIgnoreTagsConfig() *keyvaluetags.IgnoreConfig
	GetServicePackage(id string) interface{}
}

func fromMeta(meta interface{}) (*amplify.Amplify, *keyvaluetags.DefaultConfig, *keyvaluetags.IgnoreConfig) {
	m := meta.(Meta)

	return m.GetServicePackage(amplify.ServiceID).(*ServicePackage).conn, m.GetDefaultTagsConfig(), m.GetIgnoreTagsConfig()
}
