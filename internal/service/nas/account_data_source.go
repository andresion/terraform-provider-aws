package nas

import (
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aws/internal/client"
)

// See http://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/billing-getting-started.html#step-2
var billingAccountId = "386209384616"

func DataSourceBillingServiceAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAwsBillingServiceAccountRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAwsBillingServiceAccountRead(d *schema.ResourceData, meta interface{}) error {
	d.SetId(billingAccountId)
	arn := arn.ARN{
		Partition: meta.(*client.AWSClient).Partition,
		Service:   "iam",
		AccountID: billingAccountId,
		Resource:  "root",
	}.String()
	d.Set("arn", arn)

	return nil
}
