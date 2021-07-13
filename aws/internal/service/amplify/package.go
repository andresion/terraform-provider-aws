package amplify

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/amplify"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/provider/meta"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/service"
)

// Implements the service.ServicePackage interface,
type servicePackage struct {
	client *amplify.Amplify
}

func (sp *servicePackage) Client() interface{} {
	return sp.client
}

func (sp *servicePackage) EndpointsID() string {
	return amplify.EndpointsID
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

func (sp *servicePackage) Configure(ctx context.Context, sess client.ConfigProvider) error {
	sp.client = amplify.New(sess)

	return nil
}

func NewServicePackage() service.ServicePackage {
	return &servicePackage{}
}

func fromMeta(v interface{}) (*amplify.Amplify, *keyvaluetags.DefaultConfig, *keyvaluetags.IgnoreConfig) {
	m := v.(meta.Meta)

	return m.GetServicePackage(amplify.ServiceID).Client().(*amplify.Amplify), m.GetDefaultTagsConfig(), m.GetIgnoreTagsConfig()
}
