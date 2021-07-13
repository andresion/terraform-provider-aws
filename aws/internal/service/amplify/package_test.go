package amplify_test

import (
	"github.com/aws/aws-sdk-go/service/amplify"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/provider/meta"
)

func fromMeta(v interface{}) *amplify.Amplify {
	return v.(meta.Meta).GetServicePackage(amplify.ServiceID).Client().(*amplify.Amplify)
}

func resourceAwsAmplifyApp(v interface{}) *schema.Resource {
	return v.(meta.Meta).GetServicePackage(amplify.ServiceID).Resources()["aws_amplify_app"]
}

func resourceAwsAmplifyBackendEnvironment(v interface{}) *schema.Resource {
	return v.(meta.Meta).GetServicePackage(amplify.ServiceID).Resources()["aws_amplify_backend_environment"]
}

func resourceAwsAmplifyBranch(v interface{}) *schema.Resource {
	return v.(meta.Meta).GetServicePackage(amplify.ServiceID).Resources()["aws_amplify_branch"]
}

func resourceAwsAmplifyDomainAssociation(v interface{}) *schema.Resource {
	return v.(meta.Meta).GetServicePackage(amplify.ServiceID).Resources()["aws_amplify_domain_association"]
}

func resourceAwsAmplifyWebhook(v interface{}) *schema.Resource {
	return v.(meta.Meta).GetServicePackage(amplify.ServiceID).Resources()["aws_amplify_webhook"]
}
