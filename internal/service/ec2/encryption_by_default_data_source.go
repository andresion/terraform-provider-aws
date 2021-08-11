package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aws/internal/client"
)

func DataSourceEBSEncryptionByDefault() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAwsEbsEncryptionByDefaultRead,

		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}
func dataSourceAwsEbsEncryptionByDefaultRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.AWSClient).EC2Conn

	res, err := conn.GetEbsEncryptionByDefault(&ec2.GetEbsEncryptionByDefaultInput{})
	if err != nil {
		return fmt.Errorf("Error reading default EBS encryption toggle: %w", err)
	}

	d.SetId(meta.(*client.AWSClient).Region)
	d.Set("enabled", res.EbsEncryptionByDefault)

	return nil
}
